package eventproc

import (
	"log"
	"strings"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	"rapGO.io/src/bucketuploaderservice/pkg/setting"
)

func HandleVoiceFile(bi *bucket.BucketInterface, heartBeatProducer sarama.AsyncProducer, msg *sarama.ConsumerMessage) {
	var eventUUID string
	var filePrefix string
	var heartbeatDesc string
	var filenameToBucket string

	msgString := strings.Split(string(msg.Value),"/")
	filenameToBucket = msgString[len(msgString)-1]
	eventUUIDlist := strings.Split(filenameToBucket,"_")
	filePrefix = eventUUIDlist[0]
	eventUUID = strings.Split(eventUUIDlist[len(eventUUIDlist)-1],".")[0]

	switch filePrefix {
	case setting.InputPrefix():
		heartbeatDesc = "Saving raw data to cloud..."
	case setting.OutputPrefix():
		heartbeatDesc = "Saving generated data to cloud..."
	default:
		heartbeatDesc = "Internal error. File prefix not recognized."
	}

	go func(bi *bucket.BucketInterface, heartBeatProducer sarama.AsyncProducer, filenameToBucket, eventUUID, heartbeatDesc string) {
		//upload to bucket
		err := bi.Upload(filenameToBucket)
		if err != nil {
			log.Printf("The filename %s couldn't be uploaded to storage.", filenameToBucket)
			heartbeatDesc = "Internal error. The file couldn't be uploaded to the storage."
		}
		log.Printf(filenameToBucket+" has been successfully uploaded to bucket storage.")

		//Create & Emit to Kafka hearbeat topic according to content of event
		message := &sarama.ProducerMessage{Topic: setting.ToHeartbeatTopic(), Key: sarama.StringEncoder(eventUUID), Value: sarama.StringEncoder(heartbeatDesc)}
		select {
		case heartBeatProducer.Input() <- message:
			log.Println("heartbeat(@"+eventUUID+") sent : "+heartbeatDesc)
		}
	}(bi, heartBeatProducer, filenameToBucket, eventUUID, heartbeatDesc)
}