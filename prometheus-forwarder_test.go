package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
)

// var sample_input = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
// # TYPE go_gc_duration_seconds summary
// go_gc_duration_seconds{quantile="0"} 0
// go_gc_duration_seconds{quantile="0.25"} 0
// go_gc_duration_seconds{quantile="0.5"} 0
// go_gc_duration_seconds{quantile="0.75"} 0
// go_gc_duration_seconds{quantile="1"} 0
// go_gc_duration_seconds_sum 0
// go_gc_duration_seconds_count 0
// # HELP go_goroutines Number of goroutines that currently exist.
// # TYPE go_goroutines gauge
// go_goroutines 8
// # HELP go_info Information about the Go environment.
// # TYPE go_info gauge
// go_info{version="go1.13.5"} 1
// # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
// # TYPE go_memstats_alloc_bytes gauge
// go_memstats_alloc_bytes 499848
// # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
// # TYPE go_memstats_alloc_bytes_total counter
// go_memstats_alloc_bytes_total 499848
// # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
// # TYPE go_memstats_buck_hash_sys_bytes gauge
// go_memstats_buck_hash_sys_bytes 1.443544e+06
// # HELP go_memstats_frees_total Total number of frees.
// # TYPE go_memstats_frees_total counter
// go_memstats_frees_total 91
// # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
// # TYPE go_memstats_gc_cpu_fraction gauge
// go_memstats_gc_cpu_fraction 0
// # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
// # TYPE go_memstats_gc_sys_bytes gauge
// go_memstats_gc_sys_bytes 2.240512e+06
// # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
// # TYPE go_memstats_heap_alloc_bytes gauge
// go_memstats_heap_alloc_bytes 499848
// # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
// # TYPE go_memstats_heap_idle_bytes gauge
// go_memstats_heap_idle_bytes 6.4937984e+07
// # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
// # TYPE go_memstats_heap_inuse_bytes gauge
// go_memstats_heap_inuse_bytes 1.646592e+06
// # HELP go_memstats_heap_objects Number of allocated objects.
// # TYPE go_memstats_heap_objects gauge
// go_memstats_heap_objects 2504
// # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
// # TYPE go_memstats_heap_released_bytes gauge
// go_memstats_heap_released_bytes 6.4937984e+07
// # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
// # TYPE go_memstats_heap_sys_bytes gauge
// go_memstats_heap_sys_bytes 6.6584576e+07
// # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
// # TYPE go_memstats_last_gc_time_seconds gauge
// go_memstats_last_gc_time_seconds 0
// # HELP go_memstats_lookups_total Total number of pointer lookups.
// # TYPE go_memstats_lookups_total counter
// go_memstats_lookups_total 0
// # HELP go_memstats_mallocs_total Total number of mallocs.
// # TYPE go_memstats_mallocs_total counter
// go_memstats_mallocs_total 2595
// # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
// # TYPE go_memstats_mcache_inuse_bytes gauge
// go_memstats_mcache_inuse_bytes 13888
// # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
// # TYPE go_memstats_mcache_sys_bytes gauge
// go_memstats_mcache_sys_bytes 16384
// # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
// # TYPE go_memstats_mspan_inuse_bytes gauge
// go_memstats_mspan_inuse_bytes 22984
// # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
// # TYPE go_memstats_mspan_sys_bytes gauge
// go_memstats_mspan_sys_bytes 32768
// # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
// # TYPE go_memstats_next_gc_bytes gauge
// go_memstats_next_gc_bytes 4.473924e+06
// # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
// # TYPE go_memstats_other_sys_bytes gauge
// go_memstats_other_sys_bytes 1.051168e+06
// # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
// # TYPE go_memstats_stack_inuse_bytes gauge
// go_memstats_stack_inuse_bytes 524288
// # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
// # TYPE go_memstats_stack_sys_bytes gauge
// go_memstats_stack_sys_bytes 524288
// # HELP go_memstats_sys_bytes Number of bytes obtained from system.
// # TYPE go_memstats_sys_bytes gauge
// go_memstats_sys_bytes 7.189324e+07
// # HELP go_threads Number of OS threads created.
// # TYPE go_threads gauge
// go_threads 5
// # HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
// # TYPE process_cpu_seconds_total counter
// process_cpu_seconds_total 0.02
// # HELP process_max_fds Maximum number of open file descriptors.
// # TYPE process_max_fds gauge
// process_max_fds 1.048576e+06
// # HELP process_open_fds Number of open file descriptors.
// # TYPE process_open_fds gauge
// process_open_fds 7
// # HELP process_resident_memory_bytes Resident memory size in bytes.
// # TYPE process_resident_memory_bytes gauge
// process_resident_memory_bytes 7.974912e+06
// # HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
// # TYPE process_start_time_seconds gauge
// process_start_time_seconds 1.62391939735e+09
// # HELP process_virtual_memory_bytes Virtual memory size in bytes.
// # TYPE process_virtual_memory_bytes gauge
// process_virtual_memory_bytes 1.15970048e+08
// # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
// # TYPE process_virtual_memory_max_bytes gauge
// process_virtual_memory_max_bytes -1
// # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
// # TYPE promhttp_metric_handler_requests_in_flight gauge
// promhttp_metric_handler_requests_in_flight 1
// # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
// # TYPE promhttp_metric_handler_requests_total counter
// promhttp_metric_handler_requests_total{code="200"} 0
// promhttp_metric_handler_requests_total{code="500"} 0
// promhttp_metric_handler_requests_total{code="503"} 0
// # HELP uwsgi_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which uwsgi_exporter was built.
// # TYPE uwsgi_exporter_build_info gauge
// uwsgi_exporter_build_info{branch="HEAD",goversion="go1.13.5",revision="9f88775cc1a600e4038bb1eae5edfdf38f023dc4",version="1.0.0"} 1
// # HELP uwsgi_exporter_scrape_duration_seconds uwsgi_exporter: Duration of a scrape job.
// # TYPE uwsgi_exporter_scrape_duration_seconds summary
// uwsgi_exporter_scrape_duration_seconds_sum{result="success"} 0.002123745
// uwsgi_exporter_scrape_duration_seconds_count{result="success"} 1
// # HELP uwsgi_listen_queue_errors Number of listen queue errors.
// # TYPE uwsgi_listen_queue_errors gauge
// uwsgi_listen_queue_errors{stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_listen_queue_length Length of listen queue.
// # TYPE uwsgi_listen_queue_length gauge
// uwsgi_listen_queue_length{stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_signal_queue_length Length of signal queue.
// # TYPE uwsgi_signal_queue_length gauge
// uwsgi_signal_queue_length{stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_socket_can_offload Can socket offload?
// # TYPE uwsgi_socket_can_offload gauge
// uwsgi_socket_can_offload{name="0.0.0.0:7501",proto="uwsgi",stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_socket_max_queue_length Max length of socket queue.
// # TYPE uwsgi_socket_max_queue_length gauge
// uwsgi_socket_max_queue_length{name="0.0.0.0:7501",proto="uwsgi",stats_uri="http://127.0.0.1:7502/metrics"} 100
// # HELP uwsgi_socket_queue_length Length of socket queue.
// # TYPE uwsgi_socket_queue_length gauge
// uwsgi_socket_queue_length{name="0.0.0.0:7501",proto="uwsgi",stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_socket_shared Is shared socket?
// # TYPE uwsgi_socket_shared gauge
// uwsgi_socket_shared{name="0.0.0.0:7501",proto="uwsgi",stats_uri="http://127.0.0.1:7502/metrics"} 0
// # HELP uwsgi_up Whether the uwsgi server is up.
// # TYPE uwsgi_up gauge
// uwsgi_up 1
// # HELP uwsgi_worker_accepting Is this worker accepting requests?
// # TYPE uwsgi_worker_accepting gauge
// uwsgi_worker_accepting{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_accepting{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1
// # HELP uwsgi_worker_app_exceptions_total Total number of exceptions.
// # TYPE uwsgi_worker_app_exceptions_total counter
// uwsgi_worker_app_exceptions_total{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_app_exceptions_total{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_app_requests_total Total number of requests.
// # TYPE uwsgi_worker_app_requests_total counter
// uwsgi_worker_app_requests_total{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_app_requests_total{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_app_startup_time_seconds How long this app took to start.
// # TYPE uwsgi_worker_app_startup_time_seconds gauge
// uwsgi_worker_app_startup_time_seconds{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_app_startup_time_seconds{app_id="0",chdir="",mountpoint="",stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1
// # HELP uwsgi_worker_apps Number of apps.
// # TYPE uwsgi_worker_apps gauge
// uwsgi_worker_apps{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_apps{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1
// # HELP uwsgi_worker_average_response_time_seconds Average response time in seconds.
// # TYPE uwsgi_worker_average_response_time_seconds gauge
// uwsgi_worker_average_response_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0.003382
// uwsgi_worker_average_response_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_busy Is busy
// # TYPE uwsgi_worker_busy gauge
// uwsgi_worker_busy{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_busy{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_cores Number of cores.
// # TYPE uwsgi_worker_cores gauge
// uwsgi_worker_cores{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_cores{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1
// # HELP uwsgi_worker_delta_requests Number of delta requests
// # TYPE uwsgi_worker_delta_requests gauge
// uwsgi_worker_delta_requests{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_delta_requests{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_exceptions_total Total number of exceptions.
// # TYPE uwsgi_worker_exceptions_total counter
// uwsgi_worker_exceptions_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_exceptions_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_harakiri_count_total Total number of harakiri count.
// # TYPE uwsgi_worker_harakiri_count_total counter
// uwsgi_worker_harakiri_count_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_harakiri_count_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_last_spawn_time_seconds Last spawn time in seconds since epoch.
// # TYPE uwsgi_worker_last_spawn_time_seconds gauge
// uwsgi_worker_last_spawn_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1.62391939e+09
// uwsgi_worker_last_spawn_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1.62391939e+09
// # HELP uwsgi_worker_requests_total Total number of requests.
// # TYPE uwsgi_worker_requests_total counter
// uwsgi_worker_requests_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_requests_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_respawn_count_total Total number of respawn count.
// # TYPE uwsgi_worker_respawn_count_total counter
// uwsgi_worker_respawn_count_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 1
// uwsgi_worker_respawn_count_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 1
// # HELP uwsgi_worker_rss_bytes Worker RSS bytes.
// # TYPE uwsgi_worker_rss_bytes gauge
// uwsgi_worker_rss_bytes{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_rss_bytes{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_running_time_seconds Worker running time in seconds.
// # TYPE uwsgi_worker_running_time_seconds gauge
// uwsgi_worker_running_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0.006764
// uwsgi_worker_running_time_seconds{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_signal_queue_length Length of signal queue.
// # TYPE uwsgi_worker_signal_queue_length gauge
// uwsgi_worker_signal_queue_length{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_signal_queue_length{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_signals_total Total number of signals.
// # TYPE uwsgi_worker_signals_total counter
// uwsgi_worker_signals_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_signals_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_transmitted_bytes_total Worker transmitted bytes.
// # TYPE uwsgi_worker_transmitted_bytes_total counter
// uwsgi_worker_transmitted_bytes_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 109
// uwsgi_worker_transmitted_bytes_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_worker_vsz_bytes Worker VSZ bytes.
// # TYPE uwsgi_worker_vsz_bytes gauge
// uwsgi_worker_vsz_bytes{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
// uwsgi_worker_vsz_bytes{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
// # HELP uwsgi_workers Number of workers.
// # TYPE uwsgi_workers gauge
// uwsgi_workers{stats_uri="http://127.0.0.1:7502/metrics"} 2
// `

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func test_unit(t *testing.T, input, output string) {
	is_passed := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if output == buf.String() {
			is_passed = true
		} else {
			println(output)
			println(buf.String())
		}
	}
	addr := "127.0.0.1:8084"
	server := &http.Server{Addr: addr, Handler: http.HandlerFunc(handler)}
	go func() {
		server.ListenAndServe()
	}()

	var pusher = push.New(addr, "test").Format(expfmt.FmtProtoText)
	pusher, err := fill_pusher(input, pusher)
	if err != nil {
		t.Error(err)
	}
	pusher.Grouping("instance", "instance").Push()

	assertEqual(t, is_passed, true)
	server.Close()
}

func Test_fill_pusher_Gauge(t *testing.T) {
	var input_metrics = `
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 2504
`
	var desired_output = `name: "go_memstats_heap_objects"
help: ""
type: GAUGE
metric: <
  gauge: <
    value: 2504
  >
>

`
	test_unit(t, input_metrics, desired_output)
}

func Test_fill_pusher_Couner(t *testing.T) {
	var input_metrics = `
# HELP uwsgi_worker_exceptions_total Total number of exceptions.
# TYPE uwsgi_worker_exceptions_total counter
uwsgi_worker_exceptions_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="1"} 0
uwsgi_worker_exceptions_total{stats_uri="http://127.0.0.1:7502/metrics",worker_id="2"} 0
`
	var desired_output = `name: "uwsgi_worker_exceptions_total"
help: ""
type: COUNTER
metric: <
  label: <
    name: "stats_uri"
    value: "http://127.0.0.1:7502/metrics"
  >
  label: <
    name: "worker_id"
    value: "1"
  >
  counter: <
    value: 0
  >
>
metric: <
  label: <
    name: "stats_uri"
    value: "http://127.0.0.1:7502/metrics"
  >
  label: <
    name: "worker_id"
    value: "2"
  >
  counter: <
    value: 0
  >
>

`
	test_unit(t, input_metrics, desired_output)
}
