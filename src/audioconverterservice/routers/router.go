package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"rapGO.io/src/audioconverterservice/routers/api"

	"rapGO.io/src/audioconverterservice/middleware/jwt"

)

//InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/upload", api.UploadInputBLOB)

	return r
}