package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) SearchUser(username string, req_id string) (UserShortInfo, error) {
	var user UserShortInfo

	// Scan the User table for the username
	input := &dynamodb.ScanInput{
		TableName: aws.String("User"),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: username},
		},
		ProjectionExpression: aws.String("id, username"),
	}

	result, err := db.c.Scan(context.TODO(), input)
	if err != nil {
		return UserShortInfo{}, err
	}
	if len(result.Items) == 0 {
		return UserShortInfo{}, ErrUserDoesNotExist
	}

	idAttr, ok := result.Items[0]["id"].(*types.AttributeValueMemberS)
	if !ok {
		return UserShortInfo{}, fmt.Errorf("id attribute is not a string")
	}

	var banned bool
	banned, err = db.CheckBan(idAttr.Value, req_id)
	if err != nil {
		return UserShortInfo{}, err
	} else if banned {
		return UserShortInfo{}, ErrBanned
	}

	user.Id = idAttr.Value

	usernameAttr, ok := result.Items[0]["username"].(*types.AttributeValueMemberS)
	if !ok {
		return UserShortInfo{}, fmt.Errorf("username attribute is not a string")
	}
	user.Username = usernameAttr.Value

	return user, nil
}
