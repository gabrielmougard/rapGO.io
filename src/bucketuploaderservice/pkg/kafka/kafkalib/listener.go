package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/kafka/helper"
	"rapGO.io/src/bucketuploaderservice/pkg/kafka"

)

type kafkaEventListener struct {
	consumer sarama.Consumer
	partitions []int32
	mapper kafka.EventMapper
}

func NewKafkaEventListenerFromEnvironment() (kafka.EventListener, error) {
	var brokers []string
	partitions := []int32{}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	if partitionList := os.Getenv("KAFKA_PARTITIONS"); partitionList != "" {
		partitionStrings := strings.Split(partitionList, ",")
		partitions = make([]int32, len(partitionStrings))

		for i := range partitionStrings {
			partition, err := strconv.Atoi(partitionStrings[i])
			if err != nil {
				return nil, err
			}
			partitions[i] = int32(partition)
		}
	}

	client := <-helper.RetryConnect(brokers, 5*time.Second)

	return NewKafkaEventListener(client, partitions)
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (kafka.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer: consumer,
		partitions: partitions,
		mapper: kafka.NewEventMapper(),
	}

	return listener, nil
}
// interesting links : https://github.com/Shopify/sarama/issues/541
func (k *kafkaEventListener) Listen(events ...string) (<-chan kafka.Event, <-chan error, error) {
	var err error
	topics := setting.getTopicList() //several topics like `toCore`, `toBucket` and `heartbeat`
	results := make(chan kafka.Event)
	errors := make(chan error)

	partitions := k.partitions
	if len(partitions) == 0 {
		partitions, err = k.consumer.Partitions(topics)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Printf("topic %s has partitions: %v", topic)
}