package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) LikePhoto(photo_id string, user_id string) error {
	/*
		var owner_id string
		err := db.c.QueryRow("SELECT owner FROM Photo WHERE id == ?;", photo_id).Scan(&owner_id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPhotoDoesNotExist
		} else if err != nil {
			return err
		} else if user_id == owner_id {
			return ErrNotSelfLike
		}

		var banned bool
		banned, err = db.CheckBan(owner_id, user_id)
		if err != nil {
			return err
		} else if banned {
			return ErrBanned
		}

		sqlStmt := INSERT OR IGNORE INTO Like (photo, user) VALUES (?, ?)
		_, err = db.c.Exec(sqlStmt, photo_id, user_id)
		return err

	*/

	// check if the photo exists
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: photo_id},
		},
	}

	result, err := db.c.GetItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to get photo: %w", err)
	}
	if result.Item == nil {
		return ErrPhotoDoesNotExist
	}

	owner_id := result.Item["owner"].(*types.AttributeValueMemberS).Value
	if user_id == owner_id {
		return ErrNotSelfLike
	}

	// Check if user is banned
	banned, err := db.CheckBan(owner_id, user_id)
	if err != nil {
		return err
	} else if banned {
		return ErrBanned
	}

	// Add like to the likes list
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: photo_id},
		},
		UpdateExpression: aws.String("SET #likes = list_append(if_not_exists(#likes, :empty_list), :user)"),
		ExpressionAttributeNames: map[string]string{
			"#likes": "likes",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user":       &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: user_id}}},
			":empty_list": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
		},
		ConditionExpression: aws.String("attribute_not_contains(#likes, :user)"),
	}

	_, err = db.c.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		return fmt.Errorf("failed to like photo: %w", err)
	}

	return nil
}
