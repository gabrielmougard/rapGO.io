package main 

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rapGO.io/src/converterserverservice/pkg/fswatcher"
	"rapGO.io/src/converterserverservice/pkg/setting"
	"rapGO.io/src/converterserverservice/routers"
)

func init() {
	fswatcher.Setup()
}

func main() {
	gin.SetMode(setting.ServerRunMode())

	routersInit  := routers.InitRouter()
	readTimeout := 60*time.Second 
	writeTimeout := 60*time.Second
	endPoint     := fmt.Sprintf(":%d", setting.ServerHTTPport())
	maxHeaderBytes := 1 << 32

	server := &http.Server{
		Addr: endPoint,
		Handler: routersInit,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}