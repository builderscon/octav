package service

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

func defaultStorageClient(ctx context.Context) (*storage.Client, error) {
	tokesrc, err := google.DefaultTokenSource(ctx, storage.ScopeFullControl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default token source for storage client")
	}

	client, err := storage.NewClient(ctx, cloud.WithTokenSource(tokesrc))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create storage client")
	}
	return client, nil
}