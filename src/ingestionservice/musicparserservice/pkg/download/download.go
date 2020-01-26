package download

import (
	"fmt"
	"errors"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"

	"github.com/google/uuid"
	
	"rapGO.io/src/ingestionservice/musicparserservice/pkg/setting"
	
)
func downloadMP3Files(urls []string) ([][]byte, error) {
	fmt.Println("Downloading files...")
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := downloadFile(URL)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			fmt.Println("Download finished for "+URL)
			errch <- nil
		}(URL)
	}
	bytesArray := make([][]byte, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		bytesArray = append(bytesArray, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr!=""{
		err = errors.New(errStr)
	}
	return bytesArray, err
}

func downloadFile(URL string) ([]byte, error) {
	fmt.Println("Start downloading : "+URL)
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func writeFile(arr *[]byte) {
	beatPath := setting.beatPath()
	err := ioutil.WriteFile(beatPath+"beat_"+uuid.New().String()+".mp3",*arr,0644)
	if err != nil {
		panic(err)
	}
}

func writeBatch(batch *[][]byte) {
	beatPath := setting.beatPath()
	for _, song := range *batch {
		err := ioutil.WriteFile(beatPath+"beat_"+uuid.New().String()+".mp3",song,0644)
		if err != nil {
			panic(err)
		}
	}
}