package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	isOnline int = 0

	interval int
	timeout  int
	port     int
	endpoint string
)

func init() {
	flag.IntVar(&interval, "interval", 15, "Duration, in seconds, between checks")
	flag.IntVar(&timeout, "timeout", 5000, "Timeout in ms")
	flag.IntVar(&port, "port", 9091, "The port to listen on")
	flag.StringVar(&endpoint, "endpoint", "https://1.1.1.1", "The endpoint to test against")
	flag.Parse()
}

func main() {
	go (func() {
		for {
			checkOnline()
			time.Sleep(time.Minute * time.Duration(interval))
		}
	})()

	collector := newSpeedCollector()
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func checkOnline() {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	_, err := client.Get(endpoint)
	if err != nil {
		isOnline = 0
	} else {
		isOnline = 1
	}
}

type speedCollector struct {
	onlineMetric *prometheus.Desc
}

func newSpeedCollector() *speedCollector {
	return &speedCollector{
		onlineMetric: prometheus.NewDesc("is_online",
			"Is online",
			nil, nil,
		),
	}
}

func (collector *speedCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.onlineMetric
}

func (collector *speedCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.onlineMetric, prometheus.CounterValue, float64(isOnline))
}
