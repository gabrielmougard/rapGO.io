package setting

import (
	"os"
	"errors"
)

func StorageProjectID() string {
	v, ok := os.LookupEnv("STORAGE_PROJECT_ID")
	if !ok {
		panic(errors.New("the google storage project ID is not detected."))
	}
	return v
}
func StorageBucketName() string {
	v, ok := os.LookupEnv("STORAGE_BUCKET_NAME")
	if !ok {
		panic(errors.New("the google storage bucket name is not detected."))
	}
	return v
}