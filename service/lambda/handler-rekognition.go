/*
package lambdafunc

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
    sess := session.Must(session.NewSession())
    s3Client = s3.New(sess)
    rekognitionClient = rekognition.New(sess)
}

func handlerRekognition(ctx context.Context, s3Event events.S3Event) error {
    for _, record := range s3Event.Records {
        s3Record := record.S3
        bucket := s3Record.Bucket.Name
        key := s3Record.Object.Key

        fmt.Printf("Processing file: s3://%s/%s\n", bucket, key)

        // Use Rekognition to detect inappropriate content
        isSafe, err := isImageSafe(bucket, key)
        if err != nil {
            return fmt.Errorf("failed to analyze image: %v", err)
        }

        if !isSafe {
            // If the image is not safe, delete it from the bucket
            _, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(key),
            })
            if err != nil {
                return fmt.Errorf("failed to delete unsafe image: %v", err)
            }

            fmt.Printf("Unsafe image %s deleted from %s\n", key, bucket)
        }
    }

    return nil
}

func isImageSafe(bucket, key string) (bool, error) {
    result, err := rekognitionClient.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
        Image: &rekognitionTypes.Image{
            S3Object: &rekognition.S3Object{
                Bucket: aws.String(bucket),
                Name:   aws.String(key),
            },
        },
    })
    if err != nil {
        return false, err
    }

    for _, label := range result.ModerationLabels {
        if strings.Contains(*label.Name, "Explicit Nudity") || strings.Contains(*label.Name, "Violence") || strings.Contains(*label.Name, "Gore") {
            return false, nil
        }
    }

    return true, nil
}

func main() {
    lambda.Start(handler)
}
*/