package database_nosql

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"context"
	"fmt"
)

func commentToAttributeValueMap(comment Comment) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{Value: comment.Id},
        "user": &types.AttributeValueMemberS{Value: comment.User.Id},
        "created_at": &types.AttributeValueMemberS{Value: comment.CreatedDatetime},
        "text": &types.AttributeValueMemberS{Value: comment.Text},
    }
}

func (db *appdbimpl) CommentPhoto(cid string, pid string, uid string, text string, created_datetime string) (Comment, error) {
	/*
	var owner_id string
	err := db.c.QueryRow("SELECT owner FROM Photo WHERE id == ?;", pid).Scan(&owner_id)
	if errors.Is(err, sql.ErrNoRows) {
		return Comment{}, ErrPhotoDoesNotExist
	} else if err != nil {
		return Comment{}, err
	}

	var user UserShortInfo
	err = db.c.QueryRow("SELECT id, username FROM User WHERE id == ?;", uid).Scan(&user.Id, &user.Username)
	if err != nil {
		return Comment{}, err
	}

	var banned bool
	banned, err = db.CheckBan(owner_id, uid)
	if err != nil {
		return Comment{}, err
	} else if banned {
		return Comment{}, ErrBanned
	}

	comment := Comment{Id: cid, Photo: pid, User: user, CreatedDatetime: created_datetime, Text: text}

	sqlStmt := `INSERT INTO Comment (id, photo, user, text, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = db.c.Exec(sqlStmt, comment.Id, comment.Photo, comment.User.Id, comment.Text, comment.CreatedDatetime)
	return comment, err
	*/

	var owner_id string
	var user UserShortInfo
	var comment Comment

	input := &dynamodb.GetItemInput{
        TableName: aws.String("Photo"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: pid},
        },
    }

    result, err := db.c.GetItem(context.TODO(), input)
    if err != nil {
        return comment, fmt.Errorf("failed to get photo: %w", err)
    }
    if result.Item == nil {
		return comment, ErrPhotoDoesNotExist
	}

	owner_idAttr, ok := result.Item["owner"].(*types.AttributeValueMemberS)
    if !ok {
        return comment, fmt.Errorf("owner attribute is not a string")
    }
    owner_id = owner_idAttr.Value

	input = &dynamodb.GetItemInput{
        TableName: aws.String("User"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: uid},
        },
    }

    result, err = db.c.GetItem(context.TODO(), input)
    if err != nil {
        return comment, fmt.Errorf("failed to get user: %w", err)
    }
    if result.Item == nil {
		return comment, ErrUserDoesNotExist
	}

	idAttr, ok := result.Item["id"].(*types.AttributeValueMemberS)
    if !ok {
        return comment, fmt.Errorf("id attribute is not a string")
    }
    user.Id = idAttr.Value

    usernameAttr, ok := result.Item["username"].(*types.AttributeValueMemberS)
    if !ok {
        return comment, fmt.Errorf("username attribute is not a string")
    }
    user.Username = usernameAttr.Value

    var banned bool
	banned, err = db.CheckBan(owner_id, uid)
	if err != nil {
		return comment, err
	} else if banned {
		return comment, ErrBanned
	}

	comment = Comment{Id: cid, Photo: pid, User: user, CreatedDatetime: created_datetime, Text: text}
	// Update comments list
    updateInput := &dynamodb.UpdateItemInput{
	    TableName: aws.String("Photo"),
	    Key: map[string]types.AttributeValue{
	        "id": &types.AttributeValueMemberS{Value: pid},
	    },
	    UpdateExpression: aws.String("SET comments = list_append(if_not_exists(comments, :empty_list), :comment)"),
	    ExpressionAttributeValues: map[string]types.AttributeValue{
	        ":comment": &types.AttributeValueMemberL{Value: []types.AttributeValue{
	            &types.AttributeValueMemberM{Value: commentToAttributeValueMap(comment)},
	        }},
	        ":empty_list": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
	    },
	    ReturnValues: types.ReturnValueUpdatedNew,
	}

    _, err = db.c.UpdateItem(context.TODO(), updateInput)
    return comment, err
}
