package database_nosql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) UnlikePhoto(photo_id string, user_id string) error {

	// Check if photo exists
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

	// Find the index of the user_id in the likes list
	likesList := result.Item["likes"].(*types.AttributeValueMemberL).Value
	var index int
	found := false
	for i, v := range likesList {
		if v.(*types.AttributeValueMemberS).Value == user_id {
			index = i
			found = true
			break
		}
	}

	if !found {
		return ErrCannotUnlike
	}

	// Remove the user_id from the likes list
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: photo_id},
		},
		UpdateExpression: aws.String("REMOVE #likes[" + fmt.Sprintf("%d", index) + "]"),
		ExpressionAttributeNames: map[string]string{
			"#likes": "likes",
		},
		ConditionExpression: aws.String("contains(#likes, :user_id)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: user_id},
		},
	}

	_, err = db.c.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		return fmt.Errorf("failed to unlike photo: %w", err)
	}

	return nil
}
