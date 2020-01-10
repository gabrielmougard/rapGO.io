package bucket

//This is the lib for creating the interface with the google bucket

import (
	"context"
	"os"
	"io"
	"log"

	"rapGO.io/src/bucketuploaderservice/pkg/setting"
	
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

func (bi *BucketInterface) Upload(filenameToUpload string) error {
	tmpFolder := setting.TmpFolder()
	f, err := os.Open(tmpFolder+filenameToUpload)
	if err != nil {
        return err
	}
	defer f.Close()
	bucketFilename := filenameToUpload
	
	wc := bi.Client.Bucket(bi.BucketName).Object(bucketFilename).NewWriter(bi.Ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil

}