package fswatcher

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	
	fslib "rapGO.io/src/converterserverservice/pkg/fswatcher/lib"
	"rapGO.io/src/converterserverservice/pkg/setting"
)


func Setup() {
	watcher, err := fslib.NewWatcher()
	if err != nil {
	    log.Fatal(err)
	}
	producer, err := setupProducer()
	if err != nil {
		panic(err)
	}
	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	defer watcher.Close()

	done := make(chan bool)
	go func() {
	    for {
	        select {
	        case event := <-watcher.Events:
	            log.Println("event:", event)
	            if event.Op&fslib.Create == fslib.Create {
					log.Println("modified file:", event.Name)
					switch strings.Split(event.Name,"_")[0] {
					case setting.BasePrefix():
						//send to 'toBucket' and 'toCore'
						go triggerBucket(producer, signals, event.Name)
						go triggerCore(producer, signals, event.Name)
					case setting.BaseSuffix():
						//send to 'toBucket' and 'toStream'
						go triggerBucket(producer, signals, event.Name)
						go triggerStream(producer, signals, event.Name)
					default:
						log.Println("unknown prefix : "+strings.Split(event.Name,"_")[0])
					}
				}
				
	        case err, ok := <-watcher.Errors:
	            if !ok {
	                return
	            }
	            log.Println("error:", err)
	        }
	    }
	}()
	tmpFolder := setting.TmpFolder()

	err = watcher.Add(tmpFolder)
	if err != nil {
	    log.Fatal(err)
	}
	<-done
}

func setupProducer() (sarama.AsyncProducer, error) {
	kafkaBroker := setting.KafkaBroker()
	kafkaBrokers := []string{kafkaBroker}
	config := sarama.NewConfig()
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

func triggerBucket(producer sarama.AsyncProducer, signals chan os.Signal, filename string) {
	//send a kafka event to 'toBucket'
	toBucketTopic := setting.ToBucketTopic()
	//
	message := &sarama.ProducerMessage{Topic: toBucketTopic, Value: sarama.StringEncoder(filename)}
	select {
	case producer.Input() <- message:
		log.Println("event sent to "+toBucketTopic+" kafka topic.")
	case <-signals:
		producer.AsyncClose() // Trigger a shutdown of the producer.
	}
}

func triggerCore(producer sarama.AsyncProducer, signals chan os.Signal, filename string) {
	//send a kafka event to 'toCore'
	toCoreTopic := setting.toCoreTopic()
	message := &sarama.ProducerMessage{Topic: toCoreTopic, Value: sarama.StringEncoder(filename)}
	select {
	case producer.Input() <- message:
		log.Println("event sent to "+toCoreTopic+" kafka topic.")
	case <-signals:
		producer.AsyncClose() // Trigger a shutdown of the producer.
	}
}

func triggerStream(producer sarama.AsyncProducer, signals chan os.Signal, filename string) {
	//send a kafka event to 'toStream'
	toStreamTopic := setting.toCoreTopic()
	message := &sarama.ProducerMessage{Topic: toStreamTopic, Value: sarama.StringEncoder(filename)}
	select {
	case producer.Input() <- message:
		log.Println("event sent to "+toStreamTopic+" kafka topic.")
	case <-signals:
		producer.AsyncClose() // Trigger a shutdown of the producer.
	}
}
