package api

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"strings"
	"context"

	"github.com/gorilla/mux"
	"cloud.google.com/go/storage"

	"rapGO.io/src/ingestionservice/musicparserservice/pkg/download"
	"rapGO.io/src/ingestionservice/musicparserservice/pkg/setting"

)

var bucketClient *storage.Client
var bucketName string

func init() {
	//init the bucket client
	ctx = context.Background()
	projectID := setting.StorageProjectID()
	var err error
	bucketClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	bucketName = setting.StorageBucketName()
}

func handleURL(w http.ResponseWriter, r *http.Request) {
	var urls string
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	urlsSplit := strings.Split(urls,"\n")

	fmt.Println("Downloading files ...")
	batch, err := download.downloadMP3Files(urlsSplit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Download finished successfully !")
	fmt.Println("Writing results...")
	download.writeBatch(&batch)

	w.WriteHeader(200) //sucess
	return

}

func handleURLfile(w http.ResponseWriter, r *http.Request) {
	
}

func handleRawFile(w http.ResponseWriter, r *http.Request) {
	
}

func handleGenreParsing(w http.ResponseWriter, r *http.Request) {
	
}

func handleBucketData(w http.ResponseWriter, r *http.Request) {
	//create a bucket instance
	bucket := bucketClient.Bucket(bucketName)
	globalBucketMetadata := bucketstat.getBucketMetadata(&bucket) //return a json object
	objectsBucketMetadata := bucketstat.getObjectsBucketMetadata(&bucket) //return json object

	//return the merged json final object
	
}