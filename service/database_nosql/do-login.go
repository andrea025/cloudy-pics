package database_nosql

import (
	"crypto/md5"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

// Function to log in a user
func (db *appdbimpl) DoLogin(username string) (string, error) {
    input := &dynamodb.ScanInput{
        TableName: aws.String("User"),
        FilterExpression: aws.String("username = :username"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":username": &types.AttributeValueMemberS{Value: username},
        },
    }

    result, err := db.c.Scan(context.TODO(), input)
    if err != nil {
        return "", err
    }

    var id string
    if len(result.Items) == 0 {
        // If the user does not exist, generate a new user ID
        id = fmt.Sprintf("%x", md5.Sum([]byte(username)))

        // Insert the new user into the User table
        putInput := &dynamodb.PutItemInput{
            TableName: aws.String("User"),
            Item: map[string]types.AttributeValue{
                "id":       &types.AttributeValueMemberS{Value: id},
                "username": &types.AttributeValueMemberS{Value: username},
                "followers": &types.AttributeValueMemberL{
                    Value: []types.AttributeValue{}, // Empty list
                },
                "banned": &types.AttributeValueMemberL{
                    Value: []types.AttributeValue{}, // Empty list
                },
            },
        }

        _, err = db.c.PutItem(context.TODO(), putInput)
        if err != nil {
            return "", err
        }
    } else {
        // If the user exists, extract the user ID
        idAttr, ok := result.Items[0]["id"].(*types.AttributeValueMemberS)
        if !ok {
            return "", fmt.Errorf("id attribute is missing or not a string")
        }
        id = idAttr.Value
    }

    return id, nil
}
