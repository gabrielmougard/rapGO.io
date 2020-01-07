package main 

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"

	"rapGO.io/src/audioconverterservice/pkg/setting"
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