package database_nosql

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *appdbimpl) GetFollowers(id string, req_id string) ([]UserShortInfo, error) {
    users := []UserShortInfo{}

    exists, err := db.CheckUser(id)
    if err != nil {
        return nil, err
    } else if !exists {
        return nil, ErrUserDoesNotExist
    }   

    // check if the requesting user has been banned
    var banned bool
    banned, err = db.CheckBan(id, req_id)
    if err != nil {
        return nil, err
    } else if banned {
        return nil, ErrBanned
    }

    scanInput := &dynamodb.ScanInput{
        TableName: aws.String("User"),
        FilterExpression: aws.String("contains(following, :id)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":id": &types.AttributeValueMemberS{Value: id},
        },
        ProjectionExpression: aws.String("id, username"),
    }

    result, err := db.c.Scan(context.TODO(), scanInput)
    if err != nil {
        return nil, fmt.Errorf("failed to scan user table: %w", err)
    }
    if result == nil {
        return nil, fmt.Errorf("the result of GetFollowers is nil: %w", err)
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

        users = append(users, user)
    }

    return users, nil
}
