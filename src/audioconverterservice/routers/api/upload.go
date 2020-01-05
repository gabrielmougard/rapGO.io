package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rapGO.io/src/audioconverterservice/pkg/upload"

)

func UploadInputBLOB(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
}