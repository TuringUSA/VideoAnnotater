package Configs

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
)

var StoarageBucket *storage.BucketHandle
var StorageClient *storage.Client

func IntializeGcpStorage() {
	StorageClient, _ = storage.NewClient(context.TODO())
	StoarageBucket = StorageClient.Bucket(os.Getenv("GCP_STORAGE_BUCKET"))
}
