package statehandler

import (
	"strings"

	"rapGO.io/src/heartbeatservice/pkg/states"
	"rapGO.io/src/heartbeatservice/pkg/setting"

	"github.com/Shopify/sarama"
)

func HandleHeartbeat(tree *states.RbTree, msg *sarama.ConsumerMessage) {
	//get the uuid to from the key
	msgList := strings.Split(string(msg.Value),"_")
	uuid := msgList[0]
	desc := msgList[1]

	key := states.StringKey(uuid)
	//lock
	tree.Mu.Lock()
	if res, ok := tree.Get(&key); ok {
		//key exists
		if isLastHeartbeat(res) {
			//delete this node since it's useless
		} else {
			//Edit node value with the new heartbeatDesc

		}

	} else {
		//key does not exits
	}
	tree.Mu.Unlock()
}

func isLastHeartbeat(heartbeatDesc string) bool {
	if heartbeatDesc = setting.LastHeartbeatDesc() {
		return true
	} else {
		return false
	}
}