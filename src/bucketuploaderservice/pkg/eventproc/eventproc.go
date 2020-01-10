package eventproc

import (
	"log"
	"os"
	"fmt"
	"context"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/setting"

)

func ProcessEvents() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	// Start with a client
	kafkaBroker := setting.KafkaBroker()
	kafkaBrokers := []string{kafkaBroker} 
	client, err := sarama.NewClient(kafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	consumerGroupID := setting.KafkaConsumergroupID()
	toBucketTopic := setting.ToBucketTopic()

	kafkaTopics := []string{toBucketTopic} //We could listen to several topics
	defer func() { _ = client.Close() }()
	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(consumerGroupID, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()
	log.Println("Consumer up and running")
	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		handler := ConsumerGroupHandler{}

		err := group.Consume(ctx, kafkaTopics, handler)
		if err != nil {
			panic(err)
		}
	}

}


