FROM golang:1.16.5-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go get github.com/ahmetb/govvv
RUN go get -d -v ./...

RUN go build -ldflags "-X main.Version=v1.0.0"

EXPOSE 8080

# ENTRYPOINT [ "/app/prometheus-forwarder" ]