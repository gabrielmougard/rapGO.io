package bucket

//This is the lib for creating the interface with the google bucket

import (
	"context"
	"os"
	"log"
	"bytes"
	"cloud.google.com/go/storage"

)

type BucketInterface struct {
	Ctx context.Context
	ProjectID string
	Client *storage.Client
	BucketName string
	Bucket *storage.BucketHandle
}

// initialize the interface with the config values
func NewBucketInterface(projectID, bucketName string) (*BucketInterface, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}
	bucket := client.Bucket(bucketName)
	return &BucketInterface{Ctx: ctx, ProjectID: projectID, Client: client, BucketName: bucketName, Bucket: bucket}, nil
	
}

func (bi *BucketInterface) Upload(filenameToUpload string) (bool, error) {
	f, err := os.Open(config.getTmpFolder()+filenameToUpload)
	if err != nil {
        return false, err
	}
	defer f.Close()
	wc := bi.Client.Bucket(bi.Bucket).Object(object).NewWriter(bi.Ctx)
	if _, err := wc.Close(); err != nil {
		return false, err
	}
	return true, nil

}