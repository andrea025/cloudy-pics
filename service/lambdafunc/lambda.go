package lambdafunc

import (
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "errors"
    "context"
)

// AppStorage is the high level interface for the AWS Rekognition service
type AppLambda interface {
    InvokeRekognition(bucket, key string) error
    InvokeCompression(bucket, key string) error
    CheckConnectivity() error
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

// Check whether the AWS lambda is available or not (in that case, an error will be returned)
func (lambdafun *applambdaimpl) CheckConnectivity() error {
    // We can try to list bucket as a simple way to check connectivity.
    _, err := lambdafun.c.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
    return err
}
