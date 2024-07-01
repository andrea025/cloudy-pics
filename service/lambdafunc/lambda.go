package lambdafunc

import (
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "errors"
)

// AppStorage is the high level interface for the AWS Rekognition service
type AppLambda interface {
    ExecuteLambdaFunction(functionName string) error
    // Compression(filename string) error
    // CheckConnectivity() error
}

type applambdaimpl struct {
    c *lambda.Client
}

func New(lambdafun *lambda.Client) (AppLambda, error) {
	if lambdafun == nil {
		return nil, errors.New("lambda client is required when building an AppLambda")
	}

	return &applambdaimpl{
        c: lambdafun,
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