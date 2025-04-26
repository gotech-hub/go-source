package main

import (
	"context"
	"fmt"
	"go-source/pkg/metric"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const meterName = "go.opentelemetry.io/contrib/examples/prometheus"

func main() {
	ctx := context.Background()

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	// This is the equivalent of prometheus.NewHistogramVec
	histogram := metric.NewGlobalHistogramInstrument(meterName, "The latency of requests")
	_ = metric.NewHistogramWithFunc(histogram, "http", "example", func() error {
		// Do work
		return nil
	})

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
