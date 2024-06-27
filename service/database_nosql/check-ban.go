package database_nosql

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

func (db *appdbimpl) CheckBan(user_id string, target_user_id string) (bool, error) {
	/*
	var uid, tuid string
	err := db.c.QueryRow("SELECT * FROM Banned WHERE user_banning == ? AND user_banned == ?;", user_id, target_user_id).Scan(&uid, &tuid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		// considering just the error, false will be ignored
		return false, err
	}
	return true, nil
	*/

	input := &dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: user_id},
		},
		ProjectionExpression: aws.String("banned"),
	}

	result, err := db.c.GetItem(context.TODO(), input)
	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return false, fmt.Errorf("user not found")
	}

	bannedListAttr, ok := result.Item["banned"]
	if !ok {
		return false, nil // No banned list means no one is banned
	}

	bannedList, ok := bannedListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return false, fmt.Errorf("banned attribute is not a list")
	}

	for _, item := range bannedList.Value {
		bannedUserID, ok := item.(*types.AttributeValueMemberS)
		if !ok {
			continue // Skip if not a string (should not happen in a well-formed list)
		}
		if bannedUserID.Value == target_user_id {
			return true, nil
		}
	}

	return false, nil
}
