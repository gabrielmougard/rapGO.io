package uuid

import (
	"github.com/google/uuid"
	"rapGO.io/src/audioconverterservice/pkg/setting"
)

func NewVoiceUUID() string {
	basePrefix := setting.GetVoiceUUIDPrefix()
	uuid := uuid.New().String()
	baseSuffix := setting.GetVoiceUUIDSuffix()
	return basePrefix+uuid+baseSuffix //sth like input_xxxx-xxxx-xxxx-xxxx.mp3
}