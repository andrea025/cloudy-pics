package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) GetBanned(id string, req_id string) ([]UserShortInfo, error) {
	users := []UserShortInfo{}

	// Retrieve the list of banned users
	input := &dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("banned"),
	}

	result, err := db.c.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user item: %w", err)
	}

	if result.Item == nil {
		return nil, ErrUserDoesNotExist
	}

	// Check if the requesting user has been banned
	banned, err := db.CheckBan(id, req_id)
	if err != nil {
		return nil, err
	} else if banned {
		return nil, ErrBanned
	}

	bannedListAttr, ok := result.Item["banned"]
	if !ok {
		return nil, fmt.Errorf("no banned attribute in user item")
	}

	bannedList, ok := bannedListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return nil, fmt.Errorf("banned attribute is not a list")
	}

	if len(bannedList.Value) == 0 {
		return nil, nil // No banned users
	}

	keys := []map[string]types.AttributeValue{}
	for _, item := range bannedList.Value {
		if idAttr, ok := item.(*types.AttributeValueMemberS); ok {
			keys = append(keys, map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: idAttr.Value},
			})
		}
	}

	if len(keys) == 0 {
		return nil, nil // No keys to query
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			"User": {
				Keys: keys,
			},
		},
	}

	batchGetResult, err := db.c.BatchGetItem(context.TODO(), batchGetInput)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get users: %w", err)
	}

	// Print the items
	for _, item := range batchGetResult.Responses["User"] {
		var user UserShortInfo
		if idAttr, ok := item["id"].(*types.AttributeValueMemberS); ok {
			user.Id = idAttr.Value
		}
		if usernameAttr, ok := item["username"].(*types.AttributeValueMemberS); ok {
			user.Username = usernameAttr.Value
		}
		users = append(users, user)
	}

	return users, nil
}
