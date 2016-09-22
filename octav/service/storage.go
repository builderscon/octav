package service

import (
	"io"
	"os"

	"context"
)

func init() {
	switch os.Getenv("OCTAV_STORAGE_TYPE") {
	case "GoogleStorage":
		MediaStorage = &GoogleStorageClient{
			BucketName: os.Getenv("GOOGLE_STORAGE_MEDIA_BUCKET"),
		}
		CredentialStorage = &GoogleStorageClient{
			BucketName: os.Getenv("GOOGLE_STORAGE_CREDENTIAL_BUCKET"),
		}
	default:
		MediaStorage = &NullStorage{}
		CredentialStorage = &NullStorage{}
	}
}

type NullStorage struct{}
type NullObjectList struct{}

func (l NullObjectList) Next() bool {
	return false
}
func (l NullObjectList) Error() error {
	return nil
}
func (l NullObjectList) Object() interface{} {
	return nil
}

func (s *NullStorage) Move(_ context.Context, _, _ string, _ ...CallOption) error {
	return nil
}

func (c *NullStorage) Upload(_ context.Context, _ string, _ io.Reader, _ ...CallOption) error {
	return nil
}

func (c *NullStorage) Download(_ context.Context, _ string, _ io.Writer) error {
	return nil
}

func (c *NullStorage) List(_ context.Context, _ ...CallOption) (ObjectList, error) {
	return NullObjectList{}, nil
}

func (c *NullStorage) DeleteObjects(_ context.Context, _ ObjectList) error {
	return nil
}

func (C *NullStorage) URLFor(s string) string {
	return s
}
