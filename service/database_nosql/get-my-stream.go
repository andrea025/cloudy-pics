package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"context"
	"fmt"
)

func (db *appdbimpl) GetMyStream(id string) ([]Photo, error) {
	/*
	stream := []Photo{}

	sqlStmt := `SELECT id FROM Photo WHERE owner in (select user_followed from Following where user_following = ? AND user_followed not in (select user_banning from Banned where user_banned = ?)) ORDER BY created_at DESC LIMIT 50;`
	rows, err := db.c.Query(sqlStmt, id, id)
	if err != nil {
		return []Photo{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var pid string
		err = rows.Scan(&pid)
		if err != nil {
			return []Photo{}, err
		}

		var photo Photo
		photo, err = db.GetPhoto(pid, id)
		if err != nil {
			return []Photo{}, err
		}
		stream = append(stream, photo)
	}
	if err = rows.Err(); err != nil {
		return []Photo{}, err
	}

	return stream, nil
	*/

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

	var user struct {
		Following []string `dynamodbav:"following"`
	}
	err = attributevalue.UnmarshalMap(userOutput.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
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
		var bannedByUser struct {
			Id string `dynamodbav:"id"`
		}
		err = attributevalue.UnmarshalMap(item, &bannedByUser)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal banned user data: %w", err)
		}
		bannedByIds[bannedByUser.Id] = struct{}{}
	}

	// Step 3: Filter followed users to exclude users who have banned the current user
	filteredFollowingIds := []string{}
	for _, fid := range user.Following {
		if _, banned := bannedByIds[fid]; !banned {
			filteredFollowingIds = append(filteredFollowingIds, fid)
		}
	}

	// Step 4: Query photos of followed users
	for _, followedId := range filteredFollowingIds {
		photoInput := &dynamodb.QueryInput{
			TableName:              aws.String("Photo"),
			IndexName:              aws.String("owner-created_at-index"), // Assuming GSI on owner and created_at
			KeyConditionExpression: aws.String("owner = :owner"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":owner": &types.AttributeValueMemberS{Value: followedId},
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

			/*
			if len(stream) >= 50 {
				return stream, nil
			}
			*/
		}
	}

	return stream, nil
}
