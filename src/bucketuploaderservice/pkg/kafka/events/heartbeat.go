package events

import (
	"os"
)

// describe the heartbeat event to send to the `heartbeat_<UUID>` Kafka topic
type HeartbeatEvent struct {
	EventUUID     string `json:"eventUUID"`
	HeartbeatDesc string `json:"heartbeatDesc"`
}

//EventName returns the event's name
func (he *HeartbeatEvent) EventName() string {
	if prefix, ok := os.LookupEnv("KAFKA_TOPICPREFIX_HEARTBEAT"); ok {
		return prefix+he.EventUUID
	} else {
		return ""
	}
}
