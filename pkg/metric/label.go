package metric

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	InstrumentationName = "base_metric"

	HttpComponent = "http"
)

const (
	ComponentAttr = "component"
	MethodAttr    = "method"
	CodeAttr      = "code"
)

var (
	defaultTags = map[string]string{
		ComponentAttr: "",
		MethodAttr:    "",
		CodeAttr:      "SUCCESS",
	}
)

type Tags map[string]string

type LabelOption func(*Label)

type Label struct {
	tags Tags
}

func NewDefaultLabel() *Label {
	return &Label{
		tags: map[string]string{},
	}
}

func NewLabel(options ...LabelOption) *Label {
	lb := &Label{
		tags: map[string]string{},
	}

	for _, o := range options {
		o(lb)
	}

	return lb
}

func WithMethod(method string) LabelOption {
	return func(lb *Label) {
		lb.tags[MethodAttr] = method
	}
}

func WithComponent(component string) LabelOption {
	return func(lb *Label) {
		lb.tags[ComponentAttr] = component
	}
}

func WithCode(code string) LabelOption {
	return func(lb *Label) {
		lb.tags[CodeAttr] = code
	}
}

func WithAttributes(tags Tags) LabelOption {
	return func(lb *Label) {
		for k, v := range tags {
			lb.tags[k] = v
		}
	}
}

func (l *Label) GetAttributes() []attribute.KeyValue {
	tags := make([]attribute.KeyValue, 0, len(l.tags))
	for k, v := range l.tags {
		tags = append(tags, attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(v)})
	}

	return tags
}

func (l *Label) SetCode(errorCode string) {
	l.tags[CodeAttr] = errorCode
}

func NewTags(key string, value string) Tags {
	return map[string]string{key: value}
}

func NewBiTags(key1 string, value1 string, key2 string, value2 string) Tags {
	return map[string]string{key1: value1, key2: value2}
}

func NewTriTags(
	key1 string, value1 string, key2 string, value2 string, key3 string, value3 string,
) Tags {
	return map[string]string{key1: value1, key2: value2, key3: value3}
}
