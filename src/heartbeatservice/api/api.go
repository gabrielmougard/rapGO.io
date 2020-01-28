package api

import (
	"log"
	"time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"rapGO.io/src/heartbeatservice/pkg/states"
	"rapGO.io/src/heartbeatservice/pkg/setting"

)

var statesTree *states.RbTree

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Payload struct {
	Timestamp   string `json:'timestamp'`
	Description string `json:'description'`
}

func RegisterTree(tree *states.RbTree) {
	statesTree = tree
}

func SetupWebSocket(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
	}
	defer conn.Close()
	//get the coresponding node in the tree
	var wsChannel chan string

	key := states.StringKey(uuid)
	for { //wait for synchronization with statehandler
		//statesTree.Mu.Lock()
		if node, ok := statesTree.GetNode(&key); ok {
			wsChannel = node.GetDescChan()
			break
		}
		//statesTree.Mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	heartbeatCount := 0
	end := false
	for ; !end ; {
		select {
		case heartbeat := <- wsChannel:
			//write to websocket using the client
			payload := &Payload{Timestamp: time.Now().UTC().Format(time.UnixDate), Description: heartbeat}
			if err := conn.WriteJSON(payload); err != nil {
				log.Println(err)
			}
			heartbeatCount++
			if heartbeatCount == setting.TotalHeartbeatNumber() {
				statesTree.Delete(&key) //no need for the node now that all the heartbeats have been sent to the client. it also close the wsChannel.
				end = true
			}
		}
	}
}
