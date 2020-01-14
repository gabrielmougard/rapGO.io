package eventproc

import (
	"os"
	"os/signal"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/bucket"
	"rapGO.io/src/bucketuploaderservice/pkg/setting"

)

func ProcessEvents() {
	config := sarama.NewConfig()
	config.ClientID = "go-kafka-consumer"
	config.Consumer.Return.Errors = true
	kafkaBroker := setting.KafkaBroker()
  	brokers := []string{kafkaBroker}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	//create a new BucketInterface (for google storage)
	bi, err := setupBucketInterface()
	if err != nil {
		panic(err)
	}

	//create a new HeartBeatProducer
	heartbeatProducer, err := setupHeartbeatProducer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topics, _ := master.Topics()
	fmt.Println(topics)
	
	consumer, errors := consume(topics, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
				go HandleVoiceFile(bi, heartbeatProducer, msg)
			case consumerError := <-errors:
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("ProcessEvents function ended.")
}

func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}
		partitions, _ := master.Partitions(topic)
    	// this only consumes partition no 1, you would probably want to consume all partitions
		consumer, err := master.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			fmt.Printf("Topic %v Partitions: %v", topic, partitions)
			panic(err)
		}
		fmt.Println(" Start consuming topic ", topic)
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					if topic == setting.ToBucketTopic() {
						consumers <- msg
						fmt.Println("Got message on topic ", topic, msg.Value)
					} 
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}

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

func setupHeartbeatProducer() (sarama.AsyncProducer, error) {
	kafkaBroker := setting.KafkaBroker()
	kafkaBrokers := []string{kafkaBroker}
	config := sarama.NewConfig()
	sarama.Logger = log.New(os.Stderr, "[sarama_logger_heartbeat]", log.LstdFlags)
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}