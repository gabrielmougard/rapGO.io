package eventproc

import (
	"os"
	"os/signal"
	"fmt"
	//"context"
	"strings"

	"github.com/Shopify/sarama"

	"rapGO.io/src/bucketuploaderservice/pkg/setting"

)

// func ProcessEvents() {
// 	// Init config, specify appropriate version
// 	config := sarama.NewConfig()
// 	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
// 	config.Version = sarama.V2_1_0_0
// 	// Start with a client
// 	kafkaBroker := setting.KafkaBroker()
// 	kafkaBrokers := []string{kafkaBroker} 
// 	client, err := sarama.NewClient(kafkaBrokers, config)
// 	if err != nil {
// 		panic(err)
// 	}
// 	consumerGroupID := setting.KafkaConsumergroupID()
// 	toBucketTopic := setting.ToBucketTopic()

// 	kafkaTopics := []string{toBucketTopic} //We could listen to several topics
// 	defer func() { _ = client.Close() }()
// 	// Start a new consumer group

// 	//group, err := sarama.NewConsumerGroupFromClient(consumerGroupID, client)
// 	group, err := sarama.NewConsumerGroup(strings.Split(kafkaBroker,","), consumerGroupID, config)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() { _ = group.Close() }()
// 	log.Println("Consumer up and running")
// 	// Track errors
// 	go func() {
// 		for err := range group.Errors() {
// 			fmt.Println("ERROR", err)
// 		}
// 	}()
// 	// Iterate over consumer sessions.
// 	ctx := context.Background()
// 	for {
// 		handler := ConsumerGroupHandler{}

// 		err := group.Consume(ctx, kafkaTopics, handler)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// }

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
					} else {
						fmt.Println("topic not recognized : "+topic)
					}
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}