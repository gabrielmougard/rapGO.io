package eventproc

import (
	"log"
	"os"

	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	"rapGO.io/src/bucketuploaderservice/pkg/kafka"
	"rapGO.io/src/bucketuploaderservice/pkg/kafka/events"
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
	//Unmarshal content of event
	var eventUUID string
	var filePrefix string
	var heartbeatDesc string
	var filenameToBucket string
	switch e := event.(type) {
	case *eventmodel.ToBucketEvent:
		log.Printf("event %s created: %s", e.ID, e)

	default:
		log.Printf("unknown event type: %T", e)
	}
	
	go func(eh kafka.EventEmitter, filenameToBucket, eventUUID, heartbeatDesc string) {
		//upload to bucket
		err := eh.BucketInterface.Upload(filenameToBucket)
		if err != nil {
			log.Printf("The filename %s couldn't be uploaded to storage.", filenameToBucket)
			heartbeatDesc = "Internal error. The file "+filenameToBucket+" couldn't be uploaded to the storage"
		}
		//Create & Emit to Kafka hearbeat topic according to content of event
		msg := events.HeartbeatEvent{
			EventUUID: eventUUID
			HeartbeatDesc: heartbeatDesc
		}
	
		err := eh.Emit(&msg)
		if err != nil {
			log.Printf("The heartbeat for UUID: %s couldn't be sent", eventUUID)
		}
	}(ep.EventEmitter, filenameToBucket, eventUUID, heartbeatDesc)



}
