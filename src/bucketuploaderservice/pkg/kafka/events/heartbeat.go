package events

// describe the heartbeat event to send to the `heartbeat_<UUID>` Kafka topic
type HeartbeatEvent struct {
	EventUUID     string `json:"eventUUID"`
	HeartbeatDesc string `josn:"heartbeatDesc"`
}

//EventName returns the event's name
func (he *HeartbeatEvent) EventName() string {
	return "heartbeat_"+he.EventUUID
}
