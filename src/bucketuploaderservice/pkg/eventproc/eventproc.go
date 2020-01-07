package eventproc

import (
	"log"
	"os"

	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	"rapGO.io/src/bucketuploaderservice/pkg/kafka"
)
type EventProcessor struct {
	EventListener kafka.EventListener
	EventEmitter kafka.EventEmitter
	BucketInterface *bucket.BucketInterface
}

func (ep *EventProcessor) ProcessEvents() {
	log.Println("listening for events...")
	toBucketTopic, ok := os.LookupEnv("KAFKA_TOPIC_TOBUCKET")
	if !ok {
		log.Fatalf("The Kafka topic for storage bucket is not defined")
		panic()
	}
	received, errors, err := ep.EventListener.Listen(toBucketTopic)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case evt := <- received:
			fmt.Printf("got event %T: %s\n", evt, evt)
			ep.handleEvent(evt)
		case err = <-errors:
			fmt.Printf("got error while receiving event: %s\n", err)
		}
	}
}

func (ep *EventProcessor) handleEvent(event kafka.Event) {
	
}