package lambdafunc

import (
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "github.com/aws/aws-sdk-go-v2/service/lambda/types"
    "github.com/aws/aws-sdk-go-v2/aws"

    "encoding/json"
    "context"
)

var functionCompression string = "compression"

func (lambdafun *applambdaimpl) InvokeCompression(bucket, key string) error {

    payload := map[string]string{
        "bucket": bucket,
        "key":    key,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // Invoke the Lambda function
    invokeInput := &lambda.InvokeInput{
        FunctionName: aws.String(functionCompression),
        InvocationType: types.InvocationTypeRequestResponse,
        Payload: payloadBytes,
    }

    _, err = lambdafun.c.Invoke(context.TODO(), invokeInput)
    if err != nil {
        return err
    }

    return nil
}
