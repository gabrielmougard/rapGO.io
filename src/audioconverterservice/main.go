package main 

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"github.com/gin-gonic/gin"

// 	"rapGO.io/src/audioconverterservice/pkg/settings"
// 	"rapGO.io/src/audioconverterservice/routers"
// )

// func init() {
// 	setting.Setup()
// }

// func main() {
// 	gin.SetMode(setting.ServerSetting.RunMode)

// 	routersInit  := routers.InitRouter()
// 	readTimeout  := setting.ServerSetting.ReadTimeout
// 	writeTimeout := setting.ServerSetting.WriteTimeout
// 	endPoint     := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
// 	maxHeaderBytes := 1 << 20

// 	server := &http.Server{
// 		Addr: endPoint,
// 		Handler: routersInit,
// 		ReadTimeout: readTimeout,
// 		WriteTimeout: writeTimeout,
// 		MaxHeaderBytes: maxHeaderBytes,
// 	}

// 	log.Printf("[info] sstart http server listening %s", endPoint)

// 	server.ListenAndServe()
// }

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email))
	})
	router.Run(":3001")
}