package main

import (
	"fmt"
	"time"

	"rapGO.io/src/bucketuploaderservice/pkg/eventproc"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Waiting for Kafka to setup...")
	time.Sleep(60*time.Second) //security wait
	fmt.Println("eventproc setup...")

	eventproc.ProcessEvents()

}


