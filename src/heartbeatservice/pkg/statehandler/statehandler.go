package statehandler

import (
	"strings"

	"rapGO.io/src/heartbeatservice/pkg/states"

	"github.com/Shopify/sarama"
)

func HandleHeartbeat(msg *sarama.ConsumerMessage) {
	//get the uuid to from the key
	msgList := strings.Split(string(msg.Value),"_")
	uuid := msgList[0]
	desc := msgList[1]

	key := states.StringKey(uuid)

}