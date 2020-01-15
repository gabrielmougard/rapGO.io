package api

import (
	"log"
	"fmt"
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
}

type payload struct {
	timestamp   string
	description string
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

	statesTree.Mu.Lock()
	key := states.StringKey(uuid)
	if node, ok := statesTree.GetNode(&key); ok {
		wsChannel = node.GetDescChan()
	} 
	statesTree.Mu.Unlock()

	defer close(wsChannel)
	
	for {
		select {
		case heartbeat := <- wsChannel:
			//write to websocket using the client
			fmt.Println("Send heartbeat : "+heartbeat+" in webSocket")

			payload := &payload{timestamp: time.Now().UTC().Format(time.UnixDate), description: heartbeat}
			if err := conn.WriteJSON(payload); err != nil {
				log.Println(err)
			}
			log.Println("heartbeat sent successfully.")
			
			if isLastHeartbeat(heartbeat) {
				break
			}
		}
	}
}

func isLastHeartbeat(heartbeatDesc string) bool {
	if heartbeatDesc == setting.LastHeartbeatDesc() {
		return true
	} else {
		return false
	}
}