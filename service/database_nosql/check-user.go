package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "fmt"
    "context"
)

func (db *appdbimpl) CheckUser(user_id string) (bool, error) {

	input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: user_id},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return false, fmt.Errorf("failed to get user: %w", err)
    }
    if result.Item == nil {
		return false, nil
	}

	return true, nil
}
