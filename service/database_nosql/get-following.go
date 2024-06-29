package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserShortInfo struct {
	Id       string
	Username string
}


func (db *appdbimpl) GetFollowing(id string, req_id string) ([]UserShortInfo, error) {
	users := []UserShortInfo{}

	// Retrieve the list of following users
	input := &dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("following"),
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

	followingListAttr, ok := result.Item["following"]
	if !ok {
		return nil, fmt.Errorf("no following attribute in user item")
	}

	followingList, ok := followingListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return nil, fmt.Errorf("following attribute is not a list")
	}

	if len(followingList.Value) == 0 {
		return nil, nil // No followed users
	}

	keys := []map[string]types.AttributeValue{}
	for _, item := range followingList.Value {
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

	/*
	items := []map[string]types.AttributeValue{}
	err = attributevalue.UnmarshalListOfMaps(batchGetResult.Responses["User"], &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	*/

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
