package main 

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"

	"rapGO.io/src/audioconverterservice/pkg/settings"
	"rapGO.io/src/audioconverterservice/routers"
)

func init() {
	setting.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit  := routers.InitRouter()
	readTimeout  := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint     := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr: endPoint,
		Handler: routersInit,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] sstart http server listening %s", endPoint)

	server.ListenAndServe()
}