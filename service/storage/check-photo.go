package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
    "errors"
)

var ErrPhotoDoesNotExist = errors.New("photo does not exist")

func (s3_storage *apps3impl) CheckPhoto(filename string) error {
    // Check if the photo exists
    _, err := s3_storage.c.HeadObject(context.TODO(), &s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(filename),
    })

    if err != nil {
        return ErrPhotoDoesNotExist
    }

	return nil
}
