package database_nosql

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
)

type User struct {
	Id             string
	Username       string
	Followers      int
	Following      int
	UploadedPhotos int
	Photos         []Photo
}

func (db *appdbimpl) GetUserProfile(id string, req_id string) (User, error) {
    var user User

    input := &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return user, fmt.Errorf("failed to get item: %w", err)
    }

    if result.Item == nil {
        return user, ErrUserDoesNotExist
    }

    // check if the requesting user has been banned
    var banned bool
    banned, err = db.CheckBan(id, req_id)
    if err != nil {
        return user, err
    } else if banned {
        return user, ErrBanned
    }

    user.Id = id

    usernameAttr, ok := result.Item["username"].(*types.AttributeValueMemberS)
    if !ok {
        return user, fmt.Errorf("username attribute is not a string")
    }
    user.Username = usernameAttr.Value

    followingListAttr, ok := result.Item["following"]
    if !ok {
        return user, fmt.Errorf("no following attribute in user item")
    }

    followingList, ok := followingListAttr.(*types.AttributeValueMemberL)
    if !ok {
        return user, fmt.Errorf("following attribute is not a list")
    }
    user.Following = len(followingList.Value)

    followers, erro := db.GetFollowers(id, req_id)
    if erro != nil {
        return user, erro
    }
    user.Followers = len(followers)

    scanInput := &dynamodb.ScanInput{
        TableName: aws.String("Photo"),
        FilterExpression: aws.String("#owner = :owner_id"),
        ExpressionAttributeNames: map[string]string{
            "#owner": "owner",
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":owner_id": &types.AttributeValueMemberS{Value: id},
        },
        ProjectionExpression: aws.String("id"),
    }

    scanResult, err := db.c.Scan(context.TODO(), scanInput)
    if err != nil {
        return user, fmt.Errorf("failed to scan Photo table: %w", err)
    }
    
    for _, item := range scanResult.Items {
        var photo Photo

        photoIdAttr, ok := item["id"].(*types.AttributeValueMemberS)
        if !ok {
            return user, fmt.Errorf("id attribute is not a string")
        }
        photo_id := photoIdAttr.Value

        photo, err := db.GetPhoto(photo_id, id)
        if err != nil {
            return user, err
        }
        user.Photos = append(user.Photos, photo)
        user.UploadedPhotos++
    }

    return user, nil
}
