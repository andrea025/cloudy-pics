package database_nosql

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    
	"fmt"
	"context"
)

func (db *appdbimpl) GetUsers(req_id string) ([]UserShortInfo, error) {
	
	users := []UserShortInfo{}

	scanInput := &dynamodb.ScanInput{
        TableName: aws.String("User"),
        Limit:     aws.Int32(100),
    }

    result, err := db.c.Scan(context.TODO(), scanInput)
    if err != nil {
        return nil, fmt.Errorf("failed to scan users: %w", err)
    }

    if len(result.Items) == 0 {
        return nil, nil
    }

    for _, item := range result.Items {
        var user UserShortInfo

        idAttr, ok := item["id"].(*types.AttributeValueMemberS)
        if !ok {
            return nil, fmt.Errorf("id attribute is not a string")
        }
        user.Id = idAttr.Value

        usernameAttr, ok := item["username"].(*types.AttributeValueMemberS)
        if !ok {
            return nil, fmt.Errorf("username attribute is not a string")
        }
        user.Username = usernameAttr.Value

        // Check if the user is banned by the requesting user
        banned, err := db.CheckBan(user.Id, req_id)
        if err != nil {
            return nil, fmt.Errorf("failed to check ban: %w", err)
        }
        if !banned {
            users = append(users, user)
        }
    }

    return users, nil
}
