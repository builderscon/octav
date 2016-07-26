package service

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

var GoogleStorageClient *storage.Client

func init() {
	ctx := context.Background()
	tokesrc, err := google.DefaultTokenSource(ctx, storage.ScopeFullControl)
	if err != nil {
		panic("failed to get default token source for storage client: " + err.Error())
	}

	client, err := storage.NewClient(ctx, cloud.WithTokenSource(tokesrc))
	if err != nil {
		panic("failed to create storage client: " + err.Error())
	}
	GoogleStorageClient = client
}