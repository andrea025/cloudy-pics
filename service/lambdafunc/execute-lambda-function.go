package lambdafunc

import (
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "github.com/aws/aws-sdk-go-v2/service/lambda/types"
    "github.com/aws/aws-sdk-go-v2/aws"
    "context"
)

func (lambdafun *applambdaimpl) ExecuteLambdaFunction(functionName string) error {

    // Invoke the Lambda function
    invokeInput := &lambda.InvokeInput{
        FunctionName: aws.String(functionName),
        InvocationType: types.InvocationTypeRequestResponse,
    }

    _, err := lambdafun.c.Invoke(context.TODO(), invokeInput)
    if err != nil {
        return err
    }

    return nil
}
