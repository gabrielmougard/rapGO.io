package main

import (

	"errors"

	"rapGO.io/src/bucketuploaderservice/pkg/eventproc"
	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	//"rapGO.io/src/bucketuploaderservice/pkg/kafka/kafkalib"
	"rapGO.io/src/bucketuploaderservice/pkg/setting"

	"github.com/Shopify/sarama"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	//kafka config
	producer, err := setupProducer()
	if err != nil {
		panic(errors.New("Could not create AsyncProducer"))
	}

	//bucket config
	projectID := setting.StorageProjectID()
	bucketName := setting.StorageBucketName()
	pathToStorageKey := setting.StoragePath2PrivKey()

	bucketInterface, err := bucket.NewBucketInterface(projectID, bucketName)
	panicIfErr(err) 

	go eventproc.ProcessEvents()

}


