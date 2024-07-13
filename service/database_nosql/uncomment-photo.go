package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func (db *appdbimpl) UncommentPhoto(photo_id string, comment_id string, req_id string) error {

	// Fetch photo from DynamoDB
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: photo_id},
		},
		ProjectionExpression: aws.String("comments"),
	}

	result, err := db.c.GetItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to get photo: %w", err)
	}
	if result.Item == nil {
		return ErrPhotoDoesNotExist
	}

	// Get the comments list
	commentsAttr, ok := result.Item["comments"]
	if !ok {
		return fmt.Errorf("failed to get list of comments")
	}

	commentsList, ok := commentsAttr.(*types.AttributeValueMemberL)
	if !ok {
		return fmt.Errorf("comments attribute is not a list")
	}

	var index int
	var author string
	found := false

	// Iterate through the comments to find the one to delete
	for i, item := range commentsList.Value {
		itemMap, ok := item.(*types.AttributeValueMemberM)
		if !ok {
			continue
		}
		idAttr, ok := itemMap.Value["id"].(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("id attribute of comment is not a string")
		}
		idComment := idAttr.Value
		if idComment == comment_id {
			authorAttr, ok := itemMap.Value["user"].(*types.AttributeValueMemberS)
			if !ok {
				return fmt.Errorf("user attribute of comment is not a string")
			}
			author = authorAttr.Value
			if author != req_id {
				return ErrForbidden
			}

			index = i
			found = true
			break
		}
	}

	if !found {
		return ErrCommentDoesNotExist
	}

	// Update the photo to remove the comment
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: photo_id},
		},
		UpdateExpression: aws.String(fmt.Sprintf("REMOVE comments[%d]", index)),
		ReturnValues:     types.ReturnValueUpdatedNew,
	}

	_, err = db.c.UpdateItem(context.TODO(), updateInput)
	return err
}
