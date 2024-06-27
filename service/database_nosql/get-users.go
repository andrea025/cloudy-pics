package database_nosql

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    
	"fmt"
	"context"
)

func (db *appdbimpl) GetUsers(req_id string) ([]UserShortInfo, error) {
	/*
	users := []UserShortInfo{}

	sqlStmt := `SELECT * FROM User LIMIT 100;`
	rows, err := db.c.Query(sqlStmt)
	if err != nil {
		return []UserShortInfo{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var user UserShortInfo
		err = rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return []UserShortInfo{}, err
		}

		var banned bool
		banned, err = db.CheckBan(user.Id, req_id)
		if err != nil {
			return []UserShortInfo{}, err
		} else if !banned {
			users = append(users, user)
		}
	}
	if err = rows.Err(); err != nil {
		return []UserShortInfo{}, err
	}

	return users, nil
	*/

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
