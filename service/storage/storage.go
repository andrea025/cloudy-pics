package storage

import (
    "github.com/aws/aws-sdk-go-v2/service/s3"
    
    "bytes"
    "errors"
	"context"
)

var bucket string = "cloudy-pics"

// AppStorage is the high level interface for the S3 bucket
type AppStorage interface {
    UploadPhoto(filename string, image bytes.Buffer) error
    DeletePhoto(filename string) error
    CheckConnectivity() error
}

type apps3impl struct {
    c *s3.Client
}

func New(s3 *s3.Client) (AppStorage, error) {
	if s3 == nil {
		return nil, errors.New("s3 client is required when building an AppStorage")
	}

	return &apps3impl{
        c: s3,
    }, nil
}

// Check whether the s3 bucket is available or not (in that case, an error will be returned)
func (s3_storage *apps3impl) CheckConnectivity() error {
    // We can try to list bucket as a simple way to check connectivity.
    _, err := s3_storage.c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    return err
}
