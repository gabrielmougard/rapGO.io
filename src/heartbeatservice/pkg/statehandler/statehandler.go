package statehandler

import (

	"rapGO.io/src/heartbeatservice/pkg/states"
	"rapGO.io/src/heartbeatservice/pkg/setting"

	"github.com/Shopify/sarama"
)
var statesTree *states.RbTree

func RegisterTree(tree *states.RbTree) {
	statesTree = tree
}

func HandleHeartbeat(msg *sarama.ConsumerMessage) {
	//get the uuid to from the key
	
	uuid := string(msg.Key)
	desc := string(msg.Value)

	key := states.StringKey(uuid)

	//lock
	//statesTree.Mu.Lock()
	if res, ok := statesTree.Get(&key); ok {
		
		//key exists
		if isLastHeartbeat(res) {
			//delete this node since it's useless
			statesTree.EditDesc(&key, desc)
		} else {
			//Edit node value with the new heartbeatDesc
			statesTree.EditDesc(&key, desc)
		}
	} else {
		//key does not exits so insert node
		statesTree.Insert(&key, desc)
	}
	//statesTree.Mu.Unlock()
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