package main

import (
	"rapGO.io/test/testSarama/consumer"
	"rapGO.io/test/testSarama/producer"
	"os"
)

func main()  {
	if os.Args[1] == "consumer" {
		consumer.StartConsumer()
	} else if os.Args[1] == "producer" {
		producer.StartProducer()
	}
}