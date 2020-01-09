package uuid

import (
	"errors"
	"os"
	
	"github.com/google/uuid"
)

func NewVoiceUUID() string {
	// TODO : use env variables instead
	basePrefix, ok := os.LookupEnv("INPUT_PREFIX")
	if !ok {
		panic(errors.New("The environment variable INPUT_PREFIX is not defined."))
	}
	baseSuffix, ok := os.LookupEnv("INPUT_SUFFIX")
	if !ok {
		panic(errors.New("The environment variable INPUT_SUFFIX is not defined."))
	}
	uuid := uuid.New().String()
	return basePrefix+uuid+baseSuffix //sth like input_xxxx-xxxx-xxxx-xxxx.mp3
}