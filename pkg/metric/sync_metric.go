package metric

import (
	"context"
	"time"

	meter "go.opentelemetry.io/otel/metric"
)

type RecordedFunc func() error

type Metric struct {
	histogram     meter.Float64Histogram
	counter       meter.Int64Counter
	upDownCounter meter.Int64UpDownCounter

	duration time.Duration
	label    *Label
	start    time.Time
	f        RecordedFunc
}

type MetricOption func(metric *Metric)

func NewMetric(options ...MetricOption) *Metric {
	m := &Metric{}

	for _, o := range options {
		o(m)
	}

	if m.label == nil {
		m.label = NewDefaultLabel()
	}

	m.start = time.Now()
	return m
}

func (m *Metric) SetMillisDuration(duration time.Duration) *Metric {
	m.duration = duration
	return m
}

func WithHistogram(histogram meter.Float64Histogram) MetricOption {
	return func(m *Metric) {
		m.histogram = histogram
	}
}

func WithCounter(counter meter.Int64Counter) MetricOption {
	return func(m *Metric) {
		m.counter = counter
	}
}

func WithUpDownCounter(upDownCounter meter.Int64UpDownCounter) MetricOption {
	return func(m *Metric) {
		m.upDownCounter = upDownCounter
	}
}

func WithFunc(f RecordedFunc) MetricOption {
	return func(m *Metric) {
		m.f = f
	}
}

func WithLabel(options ...LabelOption) MetricOption {
	return func(m *Metric) {
		m.label = NewLabel(options...)
	}
}

func WithLabelCustomAttributes(tags map[string]string) MetricOption {
	return func(m *Metric) {
		m.label = NewLabel(WithAttributes(tags))
	}
}

func (m *Metric) Record() error {
	if m.f != nil {
		if err := m.f(); err != nil {
			m.label.SetCode(err.Error())
		} else {
			m.label.SetCode("SUCCESS")
		}
	}

	duration := time.Since(m.start).Seconds()
	if m.duration > 0 {
		duration = m.duration.Seconds()
	}
	m.histogram.Record(context.Background(), duration, meter.WithAttributes(m.label.GetAttributes()...))
	return nil
}

func (m *Metric) RecordCounter() error {
	if m.f != nil {
		if err := m.f(); err != nil {
			return err
		}
	}

	m.counter.Add(context.Background(), 1)
	return nil
}

func (m *Metric) RecordUpDownCounter() error {
	if m.f != nil {
		if err := m.f(); err != nil {
			return err
		}
	}

	m.upDownCounter.Add(context.Background(), 1)
	return nil
}

func NewHistogram(histogram meter.Float64Histogram, component, method string) *Metric {
	return NewMetric(WithLabel(WithComponent(component), WithMethod(method)), WithHistogram(histogram))
}

func NewCounter(counter meter.Int64Counter, component, method string) *Metric {
	return NewMetric(WithLabel(WithComponent(component), WithMethod(method)), WithCounter(counter))
}

func NewUpDownCounter(upDownCounter meter.Int64UpDownCounter, component, method string) *Metric {
	return NewMetric(WithLabel(WithComponent(component), WithMethod(method)), WithUpDownCounter(upDownCounter))
}

func NewHistogramWithFunc(histogram meter.Float64Histogram, component, method string, f func() error) error {
	return NewMetric(WithLabel(WithComponent(component), WithMethod(method)), WithHistogram(histogram), WithFunc(f)).Record()
}

func NewMongoDBHistogramWithFunc(component, method string, f func() error) error {
	return NewMetric(
		WithLabel(
			WithComponent(component),
			WithMethod(method),
		),
		WithHistogram(QueryMongoDBMetricHistogram),
		WithFunc(f),
	).Record()
}

func NewHttpClientHistogramDuration(component, method, code string, duration time.Duration) {
	_ = NewMetric(
		WithLabel(
			WithComponent(component),
			WithMethod(method),
			WithCode(code),
		),
		WithHistogram(HttpClientMetricHistogram),
	).SetMillisDuration(duration).Record()
}
