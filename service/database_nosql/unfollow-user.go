package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func (db *appdbimpl) UnfollowUser(user_id string, target_user_id string) error {
	/*
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `DELETE FROM Following WHERE user_following == ? AND user_followed == ?`
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
        UpdateExpression: aws.String("SET following = list_remove(following, :indices)"),
        ConditionExpression: aws.String("contains(following, :target_user_id)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":target_user_id": &types.AttributeValueMemberS{Value: target_user_id},
            ":indices": &types.AttributeValueMemberL{
                Value: []types.AttributeValue{
                    &types.AttributeValueMemberN{Value: "0"}, // Placeholder, to be calculated
                },
            },
        },
        ReturnValues: types.ReturnValueUpdatedNew,
    }

    // Find the index of the target_user_id in the following list
    followingList := result.Item["following"].(*types.AttributeValueMemberSS).Value
    var index int
    found := false
    for i, v := range followingList {
        if v == target_user_id {
            index = i
            found = true
            break
        }
    }

    if !found {
        return ErrCannotUnfollow
    }

    // Update the ExpressionAttributeValues with the correct index
    updateInput.ExpressionAttributeValues[":indices"] = &types.AttributeValueMemberL{
        Value: []types.AttributeValue{
            &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", index)},
        },
    }

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    if err != nil {
        return fmt.Errorf("failed to update following list: %w", err)
    }

    return nil
}