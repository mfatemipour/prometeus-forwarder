package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/log"
)

// will be used by govvv
var GitCommit, GitState, Version string

func handle_counter_metric(element *io_prometheus_client.MetricFamily, labels []string, pusher *push.Pusher) *push.Pusher {
	if labels == nil {
		counter := prometheus.NewCounter(prometheus.CounterOpts{Name: *element.Name})
		counter.Add(*element.Metric[0].Counter.Value)
		return pusher.Collector(counter)
	}
	counters := prometheus.NewCounterVec(prometheus.CounterOpts{Name: *element.Name}, labels)
	for _, metric := range element.Metric {
		var label_values []string
		for _, label := range metric.Label {
			label_values = append(label_values, *label.Value)
		}
		counters.WithLabelValues(label_values...).Add(*metric.Counter.Value)
	}
	return pusher.Collector(counters)
}

func handle_gauge_metric(element *io_prometheus_client.MetricFamily, labels []string, pusher *push.Pusher) *push.Pusher {
	if labels == nil {
		gauges := prometheus.NewGauge(prometheus.GaugeOpts{Name: *element.Name})
		gauges.Set(*element.Metric[0].Gauge.Value)
		return pusher.Collector(gauges)
	}
	gauges := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: *element.Name}, labels)
	for _, metric := range element.Metric {
		var label_values []string
		for _, label := range metric.Label {
			label_values = append(label_values, *label.Value)
		}

		gauges.WithLabelValues(label_values...).Set(*metric.Gauge.Value)
	}
	return pusher.Collector(gauges)
}

func scrape(pull_url string) (string, error) {
	resp, err := http.Get(pull_url)
	if err != nil {
		log.Error("scrape failed: ", err)
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

func fill_pusher(input_data string, pusher *push.Pusher) (*push.Pusher, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(strings.NewReader(input_data))
	if err != nil {
		log.Error("parse metrics failed: ", err)
		return nil, err
	}
	for _, element := range mf {
		var labels []string
		for _, label := range element.Metric[0].Label {
			labels = append(labels, *label.Name)
		}
		if *element.Type == io_prometheus_client.MetricType_COUNTER {
			pusher = handle_counter_metric(element, labels, pusher)
		} else if *element.Type == io_prometheus_client.MetricType_GAUGE {
			pusher = handle_gauge_metric(element, labels, pusher)
		} else if *element.Type == io_prometheus_client.MetricType_HISTOGRAM {
			log.Warn("histogram not supported currently")
		} else if *element.Type == io_prometheus_client.MetricType_SUMMARY {
			log.Warn("summary not supported currently")
		} else if *element.Type == io_prometheus_client.MetricType_UNTYPED {
			log.Warn("untyped not supported currently")
		}
	}
	return pusher, nil
}

func extract_and_push(input_data, push_addr, pull_url, job, instance string) error {
	var pusher = push.New(push_addr, job)
	pusher, err := fill_pusher(input_data, pusher)
	if err != nil {
		return err
	}
	err = pusher.Grouping("instance", instance).Push()
	if err != nil {
		log.Error("push failed :", err)
		return err
	} else {
		log.Info("scrape and push metrics succeded")
	}
	return nil
}

func delete_push_groups(push_addr, job, instance string) bool {
	var pusher = push.New(push_addr, job)
	var err = pusher.Grouping("instance", instance).Delete()
	if err != nil {
		log.Error("delete failed :", err)
		return false
	} else {
		log.Info("metrics deleted")
	}
	return true
}

func init_args() (int, string, []string) {
	flag.Usage = func() {
		fmt.Println("Usage: prometheus-forwarder [-i interval] -push-addr <push_addr> pull_url1,job1,instance1 pull_url2,job2,instance2 ...")
		fmt.Println("push_addr sample: 127.0.0.1:9091, pull_url sample: http://127.0.0.1/metrics")
		fmt.Println("if job and instance must be set")
		flag.PrintDefaults()
		fmt.Printf("prometheus forwarder v%s m.fatemipour@gmail.com\n", Version)
	}

	var push_addr = flag.String("push-addr", "",
		"pushgateway address 127.0.0.1:9091")
	var interval = flag.Int("interval", 5, "scrape interval in seconds")
	flag.Parse()

	var pull_urls = flag.Args()

	if len(pull_urls) == 0 || *push_addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *interval, *push_addr, pull_urls
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	interval, push_addr, pull_urls := init_args()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan bool)

	type pull_info struct {
		url      string
		job      string
		instance string
	}

	var pull_infos []pull_info

	for _, item := range pull_urls {
		s := strings.Split(item, ",")
		if len(s) != 3 {
			flag.Usage()
			panic("pull_url format error")
		}
		pull_infos = append(pull_infos, pull_info{url: s[0], job: s[1], instance: s[2]})
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, item := range pull_infos {
					if metrics_txt, err := scrape(item.url); err != nil {
						log.Error("scrape failed")
					} else if err = extract_and_push(metrics_txt, push_addr, item.url, item.job, item.instance); err != nil {
						log.Warn("push failed")
					}
				}
			case <-quit:
				return
			}
		}
	}()

	go func() {
		<-sigs
		done <- true
	}()
	<-done
	close(quit)
	ticker.Stop()
	for _, item := range pull_infos {
		delete_push_groups(push_addr, item.job, item.instance)
	}
}
