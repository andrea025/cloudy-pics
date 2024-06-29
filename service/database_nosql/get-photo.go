package database_nosql

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) GetPhoto(id string, req_id string) (Photo, error) {
	var photo Photo
	photo.Id = id

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Photo"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := db.c.GetItem(context.TODO(), input)
	if err != nil {
		return photo, fmt.Errorf("failed to get photo: %w", err)
	}
	if result.Item == nil {
		return photo, ErrPhotoDoesNotExist
	}

	owner_idAttr, ok := result.Item["owner"].(*types.AttributeValueMemberS)
	if !ok {
		return photo, fmt.Errorf("owner attribute is not a string")
	}
	owner_id := owner_idAttr.Value

	var banned bool
	banned, err = db.CheckBan(owner_id, req_id)
	if err != nil {
		return photo, err
	} else if banned {
		return photo, ErrBanned
	}
	input = &dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: owner_id},
		},
	}

	resultUser, erro := db.c.GetItem(context.TODO(), input)
	if erro != nil {
		return photo, fmt.Errorf("failed to get user: %w", err)
	}

	usernameAttr, ok := resultUser.Item["username"].(*types.AttributeValueMemberS)
	if !ok {
		return photo, fmt.Errorf("username attribute is not a string")
	}
	photo.Owner = UserShortInfo{Id: owner_id, Username: usernameAttr.Value}

	urlAttr, ok := result.Item["url"].(*types.AttributeValueMemberS)
	if !ok {
		return photo, fmt.Errorf("url attribute is not a string")
	}
	photo.PhotoUrl = urlAttr.Value

	createdDatetimeAttr, ok := result.Item["created_at"].(*types.AttributeValueMemberS)
	if !ok {
		return photo, fmt.Errorf("created_at attribute is not a string")
	}
	photo.CreatedDatetime = createdDatetimeAttr.Value

	photo.Likes = LikesCollection{Count: 0, Users: []UserShortInfo{}}
	likesListAttr, ok := result.Item["likes"]
	if !ok {
		return photo, fmt.Errorf("error in reading list of likes from photo item")
	}

	likesList, ok := likesListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return photo, fmt.Errorf("likes attribute is not a list")
	}

	if len(likesList.Value) == 0 {
		return photo, nil // No likes
	}

	// Extract user IDs from the likes list
	var userIds []string
	for _, item := range likesList.Value {
		idAttr, ok := item.(*types.AttributeValueMemberS)
		if !ok {
			return photo, fmt.Errorf("user id in the likes list is not a string")
		}
		userIds = append(userIds, idAttr.Value)
	}

	// Build keys for BatchGetItem
	likesKeys := make([]map[string]types.AttributeValue, len(userIds))
	for i, userId := range userIds {
		likesKeys[i] = map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: userId}}
	}

	if len(likesKeys) == 0 {
		return photo, nil // No keys to query
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			"User": {
				Keys: likesKeys,
				ProjectionExpression: aws.String("id, username"),
			},
		},
	}

	// Execute the BatchGetItem
	batchGetResult, err := db.c.BatchGetItem(context.TODO(), batchGetInput)
	if err != nil {
		return photo, fmt.Errorf("failed to batch get users: %w", err)
	}

	// Process the results
	users := batchGetResult.Responses["User"]
	for _, item := range users {
		var user UserShortInfo
		userIdAttr, ok := item["id"].(*types.AttributeValueMemberS)
		if !ok {
			return photo, fmt.Errorf("user id attribute is not a string")
		}
		userUsernameAttr, ok := item["username"].(*types.AttributeValueMemberS)
		if !ok {
			return photo, fmt.Errorf("user username attribute is not a string")
		}
		user.Id = userIdAttr.Value
		user.Username = userUsernameAttr.Value

		var banned bool
		banned, err = db.CheckBan(user.Id, req_id)
		if err != nil {
			return photo, err
		} else if !banned {
			photo.Likes.Users = append(photo.Likes.Users, user)
			photo.Likes.Count++
		}
	}

	photo.Comments = CommentsCollection{Count: 0, Comments: []Comment{}}
	commentsAttr, ok := result.Item["comments"]
	if !ok {
		return photo, fmt.Errorf("error in reading comments from photo item")
	}

	commentsList, ok := commentsAttr.(*types.AttributeValueMemberL)
	if !ok {
		return photo, fmt.Errorf("comments attribute is not a list")
	}

	for _, item := range commentsList.Value {
		commentAttrMap, ok := item.(*types.AttributeValueMemberM)
		if !ok {
			return photo, fmt.Errorf("comment attribute is not a map")
		}

		comment := Comment{
			Id:             commentAttrMap.Value["id"].(*types.AttributeValueMemberS).Value,
			User:           UserShortInfo{Id: commentAttrMap.Value["user"].(*types.AttributeValueMemberS).Value},
			Text:           commentAttrMap.Value["text"].(*types.AttributeValueMemberS).Value,
			CreatedDatetime: commentAttrMap.Value["created_at"].(*types.AttributeValueMemberS).Value,
		}

		var userBanned bool
		userBanned, err = db.CheckBan(comment.User.Id, req_id)
		if err != nil {
			return photo, err
		} else if !userBanned {
			input := &dynamodb.GetItemInput{
				TableName: aws.String("User"),
				Key: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: comment.User.Id},
				},
			}

			result, err := db.c.GetItem(context.TODO(), input)
			if err != nil {
				return photo, fmt.Errorf("failed to get user: %w", err)
			}

			usernameAttr, ok := result.Item["username"].(*types.AttributeValueMemberS)
			if !ok {
				return photo, fmt.Errorf("username attribute is not a string")
			}
			comment.User.Username = usernameAttr.Value

			photo.Comments.Comments = append(photo.Comments.Comments, comment)
			photo.Comments.Count++
		}
	}

	return photo, nil
}
