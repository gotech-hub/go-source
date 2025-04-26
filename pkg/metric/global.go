package metric

import meter "go.opentelemetry.io/otel/metric"

func NewGlobalHistogramInstrument(name, description string) meter.Float64Histogram {
	histogram, _ := GetGlobalMeter().Float64Histogram(
		name,
		meter.WithDescription(description),
	)
	return histogram
}

func NewGlobalCounterInstrument(name, description string) meter.Int64Counter {
	counter, _ := GetGlobalMeter().Int64Counter(
		name,
		meter.WithDescription(description),
	)
	return counter
}

func NewGlobalUpDownCounterInstrument(name, description string) meter.Int64UpDownCounter {
	counter, _ := GetGlobalMeter().Int64UpDownCounter(
		name,
		meter.WithDescription(description),
	)
	return counter
}

func NewGlobalGaugeInstrument(name, description string) meter.Int64Gauge {
	gauge, _ := GetGlobalMeter().Int64Gauge(
		name,
		meter.WithDescription(description),
	)
	return gauge
}
