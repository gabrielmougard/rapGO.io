package routers

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"rapGO.io/src/audioconverterservice/routers/api"

)

//InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/upload", api.UploadInputBLOB)

	return r
}