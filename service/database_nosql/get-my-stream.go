package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func (db *appdbimpl) GetMyStream(id string) ([]Photo, error) {
	stream := []Photo{}

	// Step 1: Get the list of users followed by the user
	userInput := &dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("following"),
	}

	userOutput, err := db.c.GetItem(context.TODO(), userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get user item: %w", err)
	}
	if userOutput.Item == nil {
		return nil, ErrUserDoesNotExist
	}

	followingListAttr, ok := userOutput.Item["following"]
	if !ok {
		return nil, fmt.Errorf("no following attribute in user item")
	}

	followingList, ok := followingListAttr.(*types.AttributeValueMemberL)
	if !ok {
		return nil, fmt.Errorf("following attribute is not a list")
	}

	// Step 2: Get the list of users who have banned the current user
	bannedByInput := &dynamodb.ScanInput{
		TableName:        aws.String("User"),
		FilterExpression: aws.String("contains(banned, :userId)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("id"),
	}

	bannedByOutput, err := db.c.Scan(context.TODO(), bannedByInput)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for banned users: %w", err)
	}

	bannedByIds := make(map[string]struct{})
	for _, item := range bannedByOutput.Items {
		idAttr, ok := item["id"].(*types.AttributeValueMemberS)
		if !ok {
			return nil, fmt.Errorf("id attribute is not a string")
		}
		bannedByIds[idAttr.Value] = struct{}{}
	}

	// Step 3: Filter followed users to exclude users who have banned the current user
	filteredFollowingIds := []string{}
	for _, fidAttr := range followingList.Value {
		fid, ok := fidAttr.(*types.AttributeValueMemberS)
		if !ok {
			return nil, fmt.Errorf("following list item is not a string")
		}
		if _, banned := bannedByIds[fid.Value]; !banned {
			filteredFollowingIds = append(filteredFollowingIds, fid.Value)
		}
	}

	// Step 4: Query photos of followed users
	for _, followedId := range filteredFollowingIds {
		photoInput := &dynamodb.QueryInput{
			TableName:              aws.String("Photo"),
			IndexName:              aws.String("owner-created_at-index"), // Assuming GSI on owner and created_at
			KeyConditionExpression: aws.String("#owner = :owner"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":owner": &types.AttributeValueMemberS{Value: followedId},
			},
			ExpressionAttributeNames: map[string]string{
            	"#owner": "owner",
        	},
			ScanIndexForward: aws.Bool(false), // Descending order
			Limit:            aws.Int32(50),
		}

		photoOutput, err := db.c.Query(context.TODO(), photoInput)
		if err != nil {
			return nil, fmt.Errorf("failed to query Photo table: %w", err)
		}

		for _, item := range photoOutput.Items {
			var photo Photo

			photoIdAttr, ok := item["id"].(*types.AttributeValueMemberS)
			if !ok {
				return nil, fmt.Errorf("id attribute is not a string")
			}
			photo_id := photoIdAttr.Value

			photo, err = db.GetPhoto(photo_id, id)
			if err != nil {
				return nil, err
			}
			stream = append(stream, photo)
		}
	}

	return stream, nil
}
