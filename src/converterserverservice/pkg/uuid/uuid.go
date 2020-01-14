package uuid

import (

	"rapGO.io/src/converterserverservice/pkg/setting"
	"github.com/google/uuid"
)

func NewVoiceUUID() string {
	inputPrefix := setting.InputPrefix()
	baseSuffix := setting.BaseSuffix()
	uuid := uuid.New().String()
	return inputPrefix+"_"+uuid+baseSuffix //sth like input_xxxx-xxxx-xxxx-xxxx.mp3
}