package main 

import (
	"fmt"
	"net/http"
	//"time"

	//"github.com/gin-gonic/gin"

	//"rapGO.io/src/converterserverservice/pkg/fswatcher"
	//"rapGO.io/src/converterserverservice/pkg/setting"
	//"rapGO.io/src/converterserverservice/routers"
)

func main() {
	fmt.Println("Waiting for Kafka to setup...")
	//time.Sleep(60*time.Second) //security wait

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

	// server.ListenAndServe()
	//fmt.Println("fswatcher setup...")
	//fswatcher.Setup()
	http.HandleFunc("/test", handler)
	http.ListenAndServe(":3001", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println("hellllllooooo !")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}