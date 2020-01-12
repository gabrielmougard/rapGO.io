package routers

import (

	"github.com/gin-gonic/gin"

	"rapGO.io/src/converterserverservice/routers/api"

)

//InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/upload", api.UploadInputBLOB)
	r.POST("/test", api.TestService)

	return r
}