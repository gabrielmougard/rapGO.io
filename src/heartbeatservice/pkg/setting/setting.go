package setting

import (
	"os"
	"strconv"
	"strings"
	"errors"
)
func KafkaBroker() string {
	v, ok := os.LookupEnv("KAFKA_BROKER")
	if !ok {
		panic(errors.New("the kafka broker is not detected."))
	}
	return v
}
func ToHeartbeatTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOHEARTBEAT")
	if !ok {
		panic(errors.New("the kafka topic toHeartbeat is not detected."))
	}
	return v
}
func ToBucketTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOBUCKET")
	if !ok {
		panic(errors.New("the kafka topic toBucket is not detected."))
	}
	return v
}
func LastHeartbeatDesc() []string {
	v, ok := os.LookupEnv("LAST_HEARTBEAT_DESC")
	if !ok {
		panic(errors.New("the last heartbeat description is not detected."))
	}
	return strings.Split(v,"|")
}
func TotalHeartbeatNumber() int {
	v, ok := os.LookupEnv("TOTAL_HEARTBEAT_NUMBER")
	if !ok {
		panic(errors.New("the last heartbeat description is not detected."))
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(errors.New("Int conversion failed for TOTAL_HEARTBEAT_NUMBER"))
	}
	return i
}