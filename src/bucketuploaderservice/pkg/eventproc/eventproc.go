package eventproc

import (
	"log"
	"os"
	"fmt"
	"strings"
	"errors"

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
		panic(errors.New("The Kafka topic for storage bucket is not defined"))
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
// So far handleEvent only receive `toBucket` events
//
func (ep *EventProcessor) handleEvent(event kafka.Event) {
	//Unmarshal content of event
	var eventUUID string
	var filePrefix string
	var heartbeatDesc string
	var filenameToBucket string

	switch e := event.(type) {
	case *events.ToBucketEvent:
		log.Printf("event %s created: %s", e.EventUUID, e)
		eventUUID = e.EventUUID
		filenameToBucket = e.EventFilename
		filePrefix = strings.Split(e.EventFilename,"_")[0]
		switch filePrefix {
		case "input":
			heartbeatDesc = "Saving raw data to cloud..."
		case "output":
			heartbeatDesc = "Saving generated data to cloud..."
		default:
			heartbeatDesc = "Internal error. File prefix not recognized."
		}
	default:
		log.Printf("unknown event type: %T", e)
	}
	
	go func(ep *EventProcessor, filenameToBucket, eventUUID, heartbeatDesc string) {
		//upload to bucket
		err := ep.BucketInterface.Upload(filenameToBucket)
		if err != nil {
			log.Printf("The filename %s couldn't be uploaded to storage.", filenameToBucket)
			heartbeatDesc = "Internal error. The file "+filenameToBucket+" couldn't be uploaded to the storage"
		}
		//Create & Emit to Kafka hearbeat topic according to content of event
		msg := events.HeartbeatEvent{
			EventUUID: eventUUID,
			HeartbeatDesc: heartbeatDesc,
		}
	
		err = ep.EventEmitter.Emit(&msg)
		if err != nil {
			log.Printf("The heartbeat for UUID: %s couldn't be sent", eventUUID)
		}
	}(ep, filenameToBucket, eventUUID, heartbeatDesc)

}
