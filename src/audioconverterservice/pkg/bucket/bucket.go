package bucket

//This is the lib for creating the interface with the google bucket

import (
	"context"
	"fmt"
	"log"
	"bytes"
	"cloud.google.com/go/storage"

	"rapGO.io/src/audioconverterservice/config"

)

type BucketInterface struct {
	Ctx context.Context
	ProjectID string
	Client *storage.Client
	BucketName string
	Bucket *storage.BucketHandle
}

// initialize the interface with the config values
func NewBucketInterface() (*BucketInterface, error) {
	ctx := context.Background()
	projectID = config.getGoogleProjectID()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}
	bucketName := config.getBucketName()
	bucket := client.Bucket(bucketName)
	return &BucketInterface{Ctx: ctx, ProjectID: projectID, Client: client, BucketName: bucketName, Bucket: bucket}, nil
	
}

func (bi *BucketInterface) Upload(filenameToUpload string) (bool, error) {

}

func (bi *BucketInterface) Dowload(filenameToDownload string) ([]bytes, error)