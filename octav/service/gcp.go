package service

import (
	"io"
	"os"
	"sync"

	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

func defaultGoogleStorageClient(ctx context.Context) (*storage.Client, error) {
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

type CallOption interface {
	Get() interface{}
}

type WithObjectAttrs storage.ObjectAttrs
type WithQueryPrefix string

func (attr WithObjectAttrs) Get() interface{} {
	return storage.ObjectAttrs(attr)
}

func (p WithQueryPrefix) Get() interface{} {
	return string(p)
}

type ObjectList interface {
	Next() bool
	Object() interface{}
	Error() error
}

type StorageClient interface {
	URLFor(string) string
	List(ctx context.Context, options ...CallOption) (ObjectList, error)
	Move(ctx context.Context, src, dst string, options ...CallOption) error
	Upload(ctx context.Context, name string, src io.Reader, options ...CallOption) error
	DeleteObjects(ctx context.Context, list ObjectList) error
}

type GoogleStorageClient struct {
	bucketOnce sync.Once
	clientOnce sync.Once
	BucketName string
	Client     *storage.Client
}

func (c *GoogleStorageClient) GetBucketName() string {
	c.bucketOnce.Do(func() {
		if c.BucketName == "" {
			c.BucketName = os.Getenv("GOOGLE_STORAGE_MEDIA_BUCKET")
		}
	})
	return c.BucketName
}

func (c *GoogleStorageClient) GetClient(ctx context.Context) *storage.Client {
	c.clientOnce.Do(func() {
		if c.Client == nil {
			client, err := defaultGoogleStorageClient(ctx)
			if err != nil {
				panic(err.Error())
			}
			c.Client = client
		}
	})
	return c.Client
}

func (c *GoogleStorageClient) URLFor(fragment string) string {
	bucketName := c.GetBucketName()
	return "http://storage.googleapis.com/" + bucketName + "/" + fragment
}

func (c *GoogleStorageClient) Move(ctx context.Context, srcName, dstName string, options ...CallOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.GoogleStorageClient.Move").BindError(&err)
		defer g.End()
	}
	storagecl := c.GetClient(ctx)
	bucketName := c.GetBucketName()
	src := storagecl.Bucket(bucketName).Object(srcName)
	dst := storagecl.Bucket(bucketName).Object(dstName)

	if pdebug.Enabled {
		pdebug.Printf("Copying %s to %s", srcName, dstName)
	}

	attrs, err := src.Attrs(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch object attrs for '%s'", srcName)
	}

	if _, err = src.CopyTo(ctx, dst, attrs); err != nil {
		return errors.Wrapf(err, "failed to copy from '%s' to '%s'", srcName, dstName)
	}

	if pdebug.Enabled {
		pdebug.Printf("Deleting %s", srcName)
	}
	if err := src.Delete(ctx); err != nil {
		return errors.Wrapf(err, "failed to delete '%s'", src)
	}
	return nil
}

func (c *GoogleStorageClient) Upload(ctx context.Context, name string, src io.Reader, options ...CallOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.GoogleStorageClient.Upload").BindError(&err)
		defer g.End()
	}
	var attrs storage.ObjectAttrs
	for _, option := range options {
		switch option.(type) {
		case WithObjectAttrs:
			attrs = option.Get().(storage.ObjectAttrs)
		}
	}

	storagecl := c.GetClient(ctx)
	wc := storagecl.Bucket(c.GetBucketName()).Object(name).NewWriter(ctx)
	wc.ObjectAttrs = attrs

	if pdebug.Enabled {
		pdebug.Printf("Writing to %s", name)
	}

	if _, err := io.Copy(wc, src); err != nil {
		return errors.Wrap(err, "failed to write data to temporary location")
	}
	// Note: DO NOT defer wc.Close(), as it's part of the write operation.
	// If wc.Close() does not complete w/o errors. the write failed
	if err := wc.Close(); err != nil {
		return errors.Wrap(err, "failed to write data to temporary location")
	}

	return nil
}

type GoogleStorageObjectList struct {
	elements <-chan interface{}
	err      error
	mu       sync.Mutex
	next     interface{}
}

func (l *GoogleStorageObjectList) Object() interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.next
}

func (l *GoogleStorageObjectList) Error() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.err
}

func (l *GoogleStorageObjectList) Next() bool {
	select {
	case next, ok := <-l.elements:
		l.mu.Lock()
		defer l.mu.Unlock()
		if !ok {
			l.elements = nil
			l.next = nil
			return false
		}

		if e, ok := next.(error); ok {
			l.err = e
			return false
		}
		l.next = next
		return true
	default:
		return false
	}
}

func (c *GoogleStorageClient) List(ctx context.Context, options ...CallOption) (ObjectList, error) {
	var q *storage.Query
	if len(options) > 0 {
		q = &storage.Query{}
	}

	for _, option := range options {
		switch option.(type) {
		case WithQueryPrefix:
			q.Prefix = option.Get().(string)
		}
	}

	out := make(chan interface{})
	go func() {
		defer close(out)
		storagecl := c.GetClient(ctx)
		b := storagecl.Bucket(c.GetBucketName())
		for q != nil {
			objects, err := b.List(ctx, q)
			if err != nil {
				return
			}

			for _, object := range objects.Results {
				out <- object
			}

			q = objects.Next
		}
	}()

	return &GoogleStorageObjectList{
		elements: out,
	}, nil
}

func (c *GoogleStorageClient) DeleteObjects(ctx context.Context, objects ObjectList) error {
	storagecl := c.GetClient(ctx)
	for objects.Next() {
		attrs, ok := objects.Object().(*storage.ObjectAttrs)
		if !ok {
			continue
		}
		if pdebug.Enabled {
			pdebug.Printf("Deleting object '%s'", attrs.Name)
		}
		if err := storagecl.Bucket(attrs.Bucket).Object(attrs.Name).Delete(ctx); err != nil {
			if pdebug.Enabled {
				pdebug.Printf("Failed to delete '%s': %s", attrs.Name, err)
			}
			return errors.Wrap(err, "failed to delete object")
		}
	}
	return nil
}
