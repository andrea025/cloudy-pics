package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) LikePhoto(photo_id string, user_id string) error {

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

	likesListAttr, ok := result.Item["likes"]
	if !ok {
		return fmt.Errorf("no likes attribute in the list")
	}

	likesList, ok := likesListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return fmt.Errorf("likes attribute is not a list")
	}

	if len(likesList.Value) != 0 {
		for _, like := range likesList.Value {
			likeS, ok := like.(*types.AttributeValueMemberS)
			if !ok {
				return fmt.Errorf("user id in the likes list is not a string")
			}

            if likeS.Value == user_id {
                // User ID already exists, no need to update
                return nil
            }
		}
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
	}

	_, err = db.c.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		return fmt.Errorf("failed to like photo: %w", err)
	}

	return nil
}
