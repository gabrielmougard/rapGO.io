package main

import (
	"time"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Shopify/sarama"

	"rapGO.io/src/heartbeatservice/api"
	"rapGO.io/src/heartbeatservice/pkg/setting"
	"rapGO.io/src/heartbeatservice/pkg/states"
	"rapGO.io/src/heartbeatservice/pkg/statehandler"

)

var statesTree *states.RbTree //global in-mem storage

func init() {
	//setup the states tree
	statesTree = states.NewRbTree()
}
func main() {
	//API server
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/heartbeat/{uuid:[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}}", api.SetupWebSocket)
		fmt.Println("Starting server on port 3002...")
		http.ListenAndServe(":3002", nil)
	}()

	//Kafka
	fmt.Println("Waiting for Kafka to setup...")
	time.Sleep(60*time.Second) //Wait for the leader election in kafka cluster

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
	consumer, errors := consumeHeartbeat(master)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
				go statehandler.HandleHeartbeat(statesTree, msg)
			case consumerError := <-errors:
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh

}

func consumeHeartbeat(master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	heartbeatTopic := setting.ToHeartbeatTopic()

	errors := make(chan *sarama.ConsumerError)
	
	//if strings.Contains(heartbeatTopic, "__consumer_offsets") {
	//	continue
	//}
	partitions, _ := master.Partitions(heartbeatTopic)
    // this only consumes partition no 1, you would probably want to consume all partitions
	consumer, err := master.ConsumePartition(heartbeatTopic, partitions[0], sarama.OffsetOldest)
	if nil != err {
		fmt.Printf("Topic %v Partitions: %v", heartbeatTopic, partitions)
		panic(err)
	}
	fmt.Println(" Start consuming topic ", heartbeatTopic)
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
	}(heartbeatTopic, consumer)
	

	return consumers, errors
}

