package events

import (
	"os"
)

// describe the core event to send to the `toCore` Kafka topic
type CoreEvent struct {
	EventUUID     string `json:"eventUUID"`
	EventFilename string `json:"eventFilename"`}

//EventName returns the event's name
func (ce *CoreEvent) EventName() string {
	if topic, ok := os.LookupEnv("KAFKA_TOPIC_TOCORE"); ok {
		return topic
	} else {
		return ""
	}
}