package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

func (db *appdbimpl) FollowUser(user_id string, target_user_id string) error {

	input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: user_id},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return err
    }
    if result.Item == nil {
		return ErrUserDoesNotExist
	}

	// Update following list
    updateInput := &dynamodb.UpdateItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: user_id},
        },
        UpdateExpression: aws.String("SET following = list_append(if_not_exists(following, :empty_list), :target_user_id)"),
        ConditionExpression: aws.String("NOT contains(following, :target_user_id)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":target_user_id": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: target_user_id}}},
            ":empty_list":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
        },
        ReturnValues: types.ReturnValueUpdatedNew,
    }

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    return err
}
