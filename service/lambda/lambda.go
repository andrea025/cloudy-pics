/*
package lambdafunc

import (
    "github.com/aws/aws-sdk-go-v2/service/rekognition"

    "errors"
	"context"
)

// AppStorage is the high level interface for the AWS Rekognition service
type AppRekognition interface {

}

type apprekoimpl struct {
    c *rekognition.Client
}

func New(reko *rekognition.Client) (AppRekognition, error) {
	if reko == nil {
		return nil, errors.New("rekognition client is required when building an AppRekognition")
	}

	return &apprekoimpl{
        c: reko,
    }, nil
}

/*
// Check whether the s3 bucket is available or not (in that case, an error will be returned)
func (reko *apprekoimpl) CheckConnectivity() error {
    // We can try to list bucket as a simple way to check connectivity.
    _, err := s3_storage.c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    return err
}
*/