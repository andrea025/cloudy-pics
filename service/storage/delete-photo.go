package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
)

func (s3_storage *apps3impl) DeletePhoto(filename string) error {

	_, err := s3_storage.c.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
	    Bucket: aws.String(bucket),
	    Key:    aws.String(filename),
	})

	return err
}
