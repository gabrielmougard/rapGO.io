	package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"rapGO.io/src/converterserverservice/pkg/uuid"

)

func UploadInputBLOB(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	}
	filename := uuid.NewVoiceUUID() //use UUID here
	fmt.Println(filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}