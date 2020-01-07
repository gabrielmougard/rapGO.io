package main

import (
	"os"
	"log"

	"rapGO.io/src/bucketuploaderservice/pkg/eventproc"
	"rapGO.io/src/bucketuploaderservice/pkg/kafka/kafkalib"
	
	"github.com/Shopify/sarama"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var eventListener kafka.EventListener
	var eventEmitter kafka.EventEmitter

	//kafka config
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	brokers, ok := os.LookupEnv("KAFKA_BROKERS")
	if !ok {
		log.Fatalf("No kafka brokers detected.")
		panic()
	}
	conn, err := sarama.NewClient([]string{brokers}, conf)
	panicIfErr(err)
	eventListener, err = kafkalib.NewKafkaEventListener(conn, []int32{})
	panicIfErr(err)
	eventEmitter, err = kafkalib.NewKafkaEventEmitter(conn)
	panicIfErr(err)

	//bucket config
	projectID, ok := os.LookupEnv("STORAGE_PROJECT_ID")
	if !ok {
		log.Fatalf("No google cloud projectID detected.")
		panic()
	}
	bucketName, ok := os.LookupEnv("STORAGE_BUCKET_NAME")
	if !ok {
		log.Fatalf("No storage bucket name detected.")
		panic()
	}
	pathToStorageKey, ok := os.LookupEnv("STORAGE_PATH_TO_BUCKET_KEY")
	if !ok {
		log.Fatalf("No pathToStorageKey to detected.")
		panic()
	} else {
		//path defined. Now, search if the file exists.
		if _, err := os.Stat(pathToStorageKey); os.IsNotExist(err) { {
			panicIfErr(err)
		}
	}

	bucketInterface, err := bucket.NewBucketInterface(projectID, bucketName)
	panicIfErr(err) 

	processor := eventproc.EventProcessor{eventListener, eventEmitter, bucketInterface}
	go processor.ProcessEvents()

}
