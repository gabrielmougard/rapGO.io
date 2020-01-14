package setting

import (
	"os"
	"errors"
	"fmt"
)

func InputPrefix() string {
	v, ok := os.LookupEnv("INPUT_PREFIX")
	if !ok {
		panic(errors.New("the input prefix is not detected."))
	}
	return v
}

func InputSuffix() string {
	v, ok := os.LookupEnv("INPUT_SUFFIX")
	if !ok {
		panic(errors.New("the input suffix is not detected."))
	}
	return v
}

func OutputPrefix() string {
	v, ok := os.LookupEnv("OUTPUT_PREFIX")
	if !ok {
		panic(errors.New("the output prefix is not detected."))
	}
	return v
}

func OutputSuffix() string {
	v, ok := os.LookupEnv("OUTPUT_SUFFIX")
	if !ok {
		panic(errors.New("the output suffix is not detected."))
	}
	return v
}

func TmpFolder() string {
	v, ok := os.LookupEnv("TMP_FOLDER")
	if !ok {
		panic(errors.New("the tmp folder is not detected."))
	}
	return v
}

func KafkaBroker() string {
	v, ok := os.LookupEnv("KAFKA_BROKER")
	if !ok {
		panic(errors.New("the kafka broker is not detected."))
	}
	return v
}

func KafkaConsumergroupID() string {
	v, ok := os.LookupEnv("KAFKA_CONSUMERGROUP_ID")
	if !ok {
		panic(errors.New("the kafka ConsumergroupID is not detected."))
	}
	return v
}

func ToBucketTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOBUCKET")
	if !ok {
		panic(errors.New("the kafka topic toBucket is not detected."))
	}
	return v
}

func ToCoreTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOCORE")
	if !ok {
		panic(errors.New("the kafka topic toCore is not detected."))
	}
	return v
}

func ToHeartbeatTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOHEARTBEAT")
	if !ok {
		panic(errors.New("the kafka topic toHeartbeat is not detected."))
	}
	return v
}

func ToStreamTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOSTREAM")
	if !ok {
		panic(errors.New("the kafka topic toStream is not detected."))
	}
	return v
}

func StorageProjectID() string {
	v, ok := os.LookupEnv("STORAGE_PROJECT_ID")
	if !ok {
		panic(errors.New("the google storage project ID is not detected."))
	}
	return v
}
func StorageBucketName() string {
	v, ok := os.LookupEnv("STORAGE_BUCKET_NAME")
	if !ok {
		panic(errors.New("the google storage bucket name is not detected."))
	}
	return v
}

func CheckStoragePath2PrivKey() {
	v, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		panic(errors.New("No pathToStorageKey to detected."))
	} else {
		//path defined. Now, search if the file exists.
		if _, err := os.Stat(v); os.IsNotExist(err) {
			panic(errors.New("The private key is not found."))
		}
		fmt.Println("The private key is detected.")
	}
}
