package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

type Photo struct {
	Id              string
	CreatedDatetime string
	PhotoUrl        string
	Owner           UserShortInfo
	Likes           LikesCollection
	Comments        CommentsCollection
}

type Comment struct {
	Id              string
	Photo           string
	User            UserShortInfo
	Text            string
	CreatedDatetime string
}

type LikesCollection struct {
	Count int
	Users []UserShortInfo
}

type CommentsCollection struct {
	Count    int
	Comments []Comment
}

func (db *appdbimpl) UploadPhoto(id string, created_at string, url string, owner string) (Photo, error) {
	
	var photo Photo
	var user UserShortInfo

	input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: owner},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return photo, fmt.Errorf("failed to get user: %w", err)
    }
    if result.Item == nil {
		return photo, ErrUserDoesNotExist
	}

	idAttr, ok := result.Item["id"].(*types.AttributeValueMemberS)
    if !ok {
        return photo, fmt.Errorf("id attribute is not a string")
    }
    user.Id = idAttr.Value

    usernameAttr, ok := result.Item["username"].(*types.AttributeValueMemberS)
    if !ok {
        return photo, fmt.Errorf("username attribute is not a string")
    }
    user.Username = usernameAttr.Value

	photo.Id, photo.CreatedDatetime, photo.PhotoUrl, photo.Owner, photo.Likes, photo.Comments = id, created_at, url, user, LikesCollection{Count: 0, Users: []UserShortInfo{}}, CommentsCollection{Count: 0, Comments: []Comment{}}	

	putInput := &dynamodb.PutItemInput{
        TableName: aws.String("Photo"),
        Item: map[string]types.AttributeValue{
            "id":       &types.AttributeValueMemberS{Value: id},
            "created_at": &types.AttributeValueMemberS{Value: created_at},
            "url":		&types.AttributeValueMemberS{Value: url},
            "owner":       &types.AttributeValueMemberS{Value: owner},
            "likes": &types.AttributeValueMemberL{
                Value: []types.AttributeValue{}, // Empty list
            },
            "comments": &types.AttributeValueMemberL{
                Value: []types.AttributeValue{}, // Empty list
            },
        },
    }

    _, err = db.c.PutItem(context.TODO(), putInput)
    if err != nil {
   	    return photo, err
    }

    return photo, nil
}
