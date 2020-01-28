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
	fmt.Println("websocket route !")
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

	for {
		statesTree.Mu.Lock()
		key := states.StringKey(uuid)
		if node, ok := statesTree.GetNode(&key); ok {
			wsChannel = node.GetDescChan()
			fmt.Println("BERLUSCONI")
			fmt.Println(wsChannel)
			break
		} else {
			fmt.Println("TREE STATE [API]")
			fmt.Println(&statesTree)
		}
		statesTree.Mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	defer close(wsChannel)
	
	for {
		select {
		case heartbeat := <- wsChannel:
			//write to websocket using the client
			fmt.Println("Send heartbeat : "+heartbeat+" in webSocket")

			payload := &Payload{Timestamp: time.Now().UTC().Format(time.UnixDate), Description: heartbeat}
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

func TestServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{\"hello\": \"world\"}"))
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