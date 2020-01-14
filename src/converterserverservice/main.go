package main 

import (
	//"fmt"
	"net/http"
	"github.com/rs/cors"
	"time"
	"fmt"

	//"github.com/gin-gonic/gin"

	"rapGO.io/src/converterserverservice/pkg/fswatcher"
	//"rapGO.io/src/converterserverservice/pkg/setting"
	//"rapGO.io/src/converterserverservice/routers"
	"rapGO.io/src/converterserverservice/routers/api"

)

func main() {
	fmt.Println("Waiting for Kafka to setup...")
	time.Sleep(60*time.Second) //security wait

	// gin.SetMode(setting.ServerRunMode())

	// routersInit  := routers.InitRouter()
	// readTimeout := 60*time.Second 
	// writeTimeout := 60*time.Second
	// endPoint     := fmt.Sprintf(":%d", setting.ServerHTTPport())
	// maxHeaderBytes := 1 << 32

	// server := &http.Server{
	// 	Addr: endPoint,
	// 	Handler: routersInit,
	// 	ReadTimeout: readTimeout,
	// 	WriteTimeout: writeTimeout,
	// 	MaxHeaderBytes: maxHeaderBytes,
	// }

	// fmt.Println("[info] start http server listening %s", endPoint)

	//server.ListenAndServe()
	//fmt.Println("fswatcher setup...")
	go fswatcher.Setup()
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", api.UploadInputBLOB)
	handler := cors.Default().Handler(mux)
	fmt.Println("Starting server on port 3001...")
	http.ListenAndServe(":3001", handler)
}

