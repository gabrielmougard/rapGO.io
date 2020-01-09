package events

import (
	"os"
)

type ToBucketEvent struct {
	EventUUID     string `json:"eventUUID"`
	EventFilename string `json:"eventFilename"`
}

func (be *ToBucketEvent) EventName() string {
	if topic, ok := os.LookupEnv("KAFKA_TOPIC_TOBUCKET"); ok {
		return topic
	} else {
		return ""
	}
}