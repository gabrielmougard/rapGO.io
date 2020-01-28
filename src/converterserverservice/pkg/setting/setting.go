package setting

import (
	"os"
	"errors"
)

func ServerHTTPport() string {
	v, ok := os.LookupEnv("SERVER_HTTP_PORT")
	if !ok {
		panic(errors.New("the server port is not detected."))
	}
	return v
}

func ServerRunMode() string {
	v, ok := os.LookupEnv("SERVER_RUN_MODE")
	if !ok {
		panic(errors.New("the server run mode is not detected."))
	}
	return v
}

func ServerReadTimeout() string {
	v, ok := os.LookupEnv("SERVER_READ_TIMEOUT")
	if !ok {
		panic(errors.New("the server ReadTimeout is not detected."))
	}
	return v
}

func ServerWriteTimeout() string {
	v, ok := os.LookupEnv("SERVER_WRITE_TIMEOUT")
	if !ok {
		panic(errors.New("the server WriteTimeout is not detected."))
	}
	return v
}

func InputPrefix() string {
	v, ok := os.LookupEnv("INPUT_PREFIX")
	if !ok {
		panic(errors.New("the input prefix is not detected."))
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

func BaseSuffix() string {
	v, ok := os.LookupEnv("INPUT_SUFFIX")
	if !ok {
		panic(errors.New("the input suffix is not detected."))
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

func ToStreamTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_TOSTREAM")
	if !ok {
		panic(errors.New("the kafka topic toStream is not detected."))
	}
	return v
}

func ToHeartbeatTopic() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_HEARTBEAT")
	if !ok {
		panic(errors.New("the kafka topic toHeartbeat is not detected."))
	}
	return v
}