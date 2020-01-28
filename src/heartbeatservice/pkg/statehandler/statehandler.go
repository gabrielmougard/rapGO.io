package statehandler

import (

	"rapGO.io/src/heartbeatservice/pkg/states"
	"rapGO.io/src/heartbeatservice/pkg/setting"

	"github.com/Shopify/sarama"
)

func HandleHeartbeat(tree *states.RbTree, msg *sarama.ConsumerMessage) {
	//get the uuid to from the key

	uuid := string(msg.Key)
	desc := string(msg.Value)

	key := states.StringKey(uuid)
	//lock
	tree.Mu.Lock()
	if res, ok := tree.Get(&key); ok {
		//key exists
		if isLastHeartbeat(res) {
			//delete this node since it's useless
			tree.Delete(&key)
		} else {
			//Edit node value with the new heartbeatDesc
			tree.EditDesc(&key, desc)
		}
	} else {
		//key does not exits so insert node
		tree.Insert(&key, desc)
	}
	tree.Mu.Unlock()
}

func isLastHeartbeat(heartbeatDesc string) bool {
	possibleHeartbeats := setting.LastHeartbeatDesc()

	if heartbeatDesc == possibleHeartbeats[0] {
		return true
	} else if heartbeatDesc ==  possibleHeartbeats[1] {
		return true
	} else {
		return false
	}
}