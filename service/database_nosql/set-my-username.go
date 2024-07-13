package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

func (db *appdbimpl) SetMyUsername(id string, username string) error {
	
	input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return err
    }
    if result.Item == nil {
		return ErrUserDoesNotExist
	}

	// Update username
    updateInput := &dynamodb.UpdateItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
        UpdateExpression: aws.String("SET username = :username"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":username": &types.AttributeValueMemberS{Value: username},
        },
        ReturnValues: types.ReturnValueUpdatedNew,
    }

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    return err
}
