package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func (db *appdbimpl) UnbanUser(user_id string, target_user_id string) error {
	/*
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `DELETE FROM Banned WHERE user_banning == ? AND user_banned == ?`
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

	// Update banned list
    updateInput := &dynamodb.UpdateItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: user_id},
        },
        UpdateExpression: aws.String("SET banned = list_remove(banned, :indices)"),
        ConditionExpression: aws.String("contains(banned, :target_user_id)"),
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

    // Find the index of the target_user_id in the banned list
    bannedList := result.Item["banned"].(*types.AttributeValueMemberSS).Value
    var index int
    found := false
    for i, v := range bannedList {
        if v == target_user_id {
            index = i
            found = true
            break
        }
    }

    if !found {
        return ErrCannotUnban
    }

    // Update the ExpressionAttributeValues with the correct index
    updateInput.ExpressionAttributeValues[":indices"] = &types.AttributeValueMemberL{
        Value: []types.AttributeValue{
            &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", index)},
        },
    }

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    if err != nil {
        return fmt.Errorf("failed to update banned list: %w", err)
    }

    return nil
}
