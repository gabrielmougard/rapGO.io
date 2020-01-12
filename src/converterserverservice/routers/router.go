package routers

import (

	"github.com/gin-gonic/gin"

	"rapGO.io/src/converterserverservice/routers/api"

)

//InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	r.Use(Default())

	r.POST("/upload", api.UploadInputBLOB)

	return r
}