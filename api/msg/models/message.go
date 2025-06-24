package models

type ConsumerMessage struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}
