package eventproc

import (
	"fmt"
	"os"
	"log"
	"strings"
	"os/signal"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	"rapGO.io/src/bucketuploaderservice/pkg/setting"
)

func setupBucketInterface() (*bucket.BucketInterface, error) {
	//bucket config
	projectID := setting.StorageProjectID()
	bucketName := setting.StorageBucketName()
	setting.CheckStoragePath2PrivKey()

	bucketInterface, err := bucket.NewBucketInterface(projectID, bucketName)
	if err != nil {
		panic(err)
	}
	return bucketInterface, nil
}

func setupProducer() (sarama.AsyncProducer, error) {

	kafkaBroker := setting.KafkaBroker()
	kafkaBrokers := []string{kafkaBroker}
	config := sarama.NewConfig()
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

// ConsumerGroupHandler represents the sarama consumer group
type ConsumerGroupHandler struct{}

// Setup is run before consumer start consuming, is normally used to setup things such as database connections
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages(), here is supposed to be what you want to
// do with the message. In this example the message will be logged with the topic name, partition and message value.
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var eventUUID string
	var filePrefix string
	var heartbeatDesc string
	var filenameToBucket string

	producer, err := setupProducer()
	if err != nil {
		panic(err)
	}
	toBucketTopic := setting.ToBucketTopic()
	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for msg := range claim.Messages() {

		switch msg.Topic {
		case toBucketTopic:
			filenameToBucket = string(msg.Value)
			exploded := strings.Split(filenameToBucket,"_")
			filePrefix = exploded[0]
			eventUUID = strings.Split(exploded[1],".")[0]

			switch filePrefix {
			case setting.InputPrefix():
				heartbeatDesc = "Saving raw data to cloud..."
			case setting.OutputPrefix():
				heartbeatDesc = "Saving generated data to cloud..."
			default:
				heartbeatDesc = "Internal error. File prefix not recognized."
			}
			//create the bucket interface
			bucketInterface, err := setupBucketInterface()
			if err != nil {
				panic(err)
			}

			go func(bi *bucket.BucketInterface, producer sarama.AsyncProducer,  signals chan os.Signal, filenameToBucket, eventUUID, heartbeatDesc string) {
				//upload to bucket
				err := bi.Upload(filenameToBucket)
				if err != nil {
					log.Printf("The filename %s couldn't be uploaded to storage.", filenameToBucket)
					heartbeatDesc = "Internal error. The file "+filenameToBucket+" couldn't be uploaded to the storage"
				}
				//Create & Emit to Kafka hearbeat topic according to content of event
				message := &sarama.ProducerMessage{Topic: setting.ToHeartbeatTopicPrefix()+eventUUID, Value: sarama.StringEncoder(heartbeatDesc)}
				select {
				case producer.Input() <- message:
					log.Println("heartbeat sent : "+heartbeatDesc)
				case <-signals:
					producer.AsyncClose() // Trigger a shutdown of the producer.
					return
				}
				
			}(bucketInterface, producer, signals, filenameToBucket, eventUUID, heartbeatDesc)

		default:
			fmt.Printf("The topic is not recognized.")
		}
	}
	return nil
}
