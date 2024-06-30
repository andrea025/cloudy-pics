package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
	"bytes"
)

func (s3_storage *apps3impl) UploadPhoto(filename string, image bytes.Buffer) error {

	_, err := s3_storage.c.PutObject(context.TODO(), &s3.PutObjectInput{
	  	Bucket: aws.String(bucket),
	  	Key:    aws.String(filename),
	  	Body:   bytes.NewReader(image.Bytes()),
	})

	// WAIT A BIT.............
	// IF getObject DOES NOT RETURN THE PIC, GENERATE A CUSTOM ERORR THAT SAYS "AYO! PORNO!"


	return err
}
