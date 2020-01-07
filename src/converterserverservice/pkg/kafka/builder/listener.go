package builder

import (
	"errors"
	"log"
	"os"
	"rapGO.io/src/converterserverservice/pkg/kafka"
	"rapGO.io/src/converterserverservice/pkg/kafka/lib"
)

func NewEventListenerFromEnvironment() (kafka.EventListener, error) {
	var listener kafka.EventListener
	var err error
	
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)
		listener, err = lib.NewKafkaEventListenerFromEnvironment()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Neither KAFKA_BROKERS or another brokers(unsupported for now) specified")
	}

	return listener, nil
}