package event_bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type ErrCallback chan error
type Subscriber chan *Event
type HandleEvent func(event *Event) error

type Event struct {
	Ctx         context.Context
	Data        interface{}
	ErrCallback ErrCallback
	Response    interface{}
}

type EventBus struct {
	subscriber map[string]Subscriber
	mu         sync.RWMutex
}

var (
	eventBus *EventBus
	once     sync.Once
)

func NewEventBus() *EventBus {
	once.Do(func() {
		eventBus = &EventBus{
			subscriber: make(map[string]Subscriber),
		}
	})
	return eventBus
}

func GetEventBus() *EventBus {
	return eventBus
}

func (eb *EventBus) Subscribe(eventType string, fn HandleEvent) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscriber[eventType]; !ok {
		eb.subscriber[eventType] = make(Subscriber, 10)
	}

	sub := eb.subscriber[eventType]

	go func(sub Subscriber) {
		for event := range sub {
			select {
			case <-event.Ctx.Done():
			default:
				sendErr(event.ErrCallback, fn(event))
			}
		}
	}(sub)
}

func (eb *EventBus) Publish(eventType string, event *Event) error {
	if event.ErrCallback == nil {
		event.ErrCallback = make(ErrCallback)
	}

	if event.Response != nil {
		v := reflect.ValueOf(event.Response)
		if v.Kind() != reflect.Ptr {
			close(event.ErrCallback)
			return fmt.Errorf("response must be a pointer")
		}
	}

	if err := eb.publish(eventType, event); err != nil {
		return err
	}

	defer close(event.ErrCallback)

	// Wait for the ErrCallback
	select {
	case <-event.Ctx.Done():
		return fmt.Errorf("%w: no ErrCallback received", event.Ctx.Err())
	case err := <-event.ErrCallback:
		return err
	}
}

func (eb *EventBus) publish(eventType string, event *Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	sub, ok := eb.subscriber[eventType]
	if !ok {
		return fmt.Errorf("no subscriber found for event type %s", eventType)
	}

	// Publish event to subscriber
	sub <- event

	return nil
}

// Safe function to send errors into ErrCallback, checks if channel is open
func sendErr(ch ErrCallback, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	select {
	case ch <- err:
	default:
		// If channel is already closed or no receivers, ignore
	}
}
