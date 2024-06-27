package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

func (db *appdbimpl) FollowUser(user_id string, target_user_id string) error {
	/*
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `INSERT OR IGNORE INTO Following (user_following, user_followed) VALUES (?, ?)`
	_, err = db.c.Exec(sqlStmt, user_id, target_user_id)
	return err
	*/

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
            ":target_user_id": &types.AttributeValueMemberSS{Value: []string{target_user_id}},
            ":empty_list":     &types.AttributeValueMemberSS{Value: []string{}},
        },
        ReturnValues: types.ReturnValueUpdatedNew,
    }

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    return err
}
