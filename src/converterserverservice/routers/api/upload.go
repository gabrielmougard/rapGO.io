package api

import (
	"net/http"
	"fmt"
	
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
	// TODO :
	// create a new kafka topic here for the heartbeat. the topic name should be : heartbeat_<UUID>
	// how to : https://stackoverflow.com/questions/44094926/creating-kafka-topic-in-sarama
	//
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}

func TestService( c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Hello gin server !"))
}