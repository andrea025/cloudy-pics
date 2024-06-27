package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func (db *appdbimpl) SearchUser(username string, req_id string) (UserShortInfo, error) {
	/*
	var user UserShortInfo
	err := db.c.QueryRow("SELECT * FROM User WHERE username == ?;", username).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserDoesNotExist
	} else if err != nil {
		return user, err
	}

	var banned bool
	banned, err = db.CheckBan(user.Id, req_id)
	if err != nil {
		return user, err
	} else if banned {
		return UserShortInfo{}, nil
	}

	return user, nil
	*/

	var user UserShortInfo
	// Query the User table for the username
    input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "username": &types.AttributeValueMemberS{Value: username},
        },
        ProjectionExpression: aws.String("id, username"),
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return UserShortInfo{}, err
    }
    if result.Item == nil {
    	return UserShortInfo{}, ErrUserDoesNotExist
    }

    idAttr, ok := result.Item["id"].(*types.AttributeValueMemberS)
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

    usernameAttr, ok := result.Item["username"].(*types.AttributeValueMemberS)
    if !ok {
        return UserShortInfo{}, fmt.Errorf("username attribute is not a string")
    }
    user.Username = usernameAttr.Value


    return user, nil
}
