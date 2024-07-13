package database_nosql

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "context"
    "fmt"
)

func (db *appdbimpl) DeletePhoto(photo_id string, req_id string) (string, error) {

    getItemInput := &dynamodb.GetItemInput{
        TableName: aws.String("Photo"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: photo_id},
        },
    }

    getItemOutput, err := db.c.GetItem(context.TODO(), getItemInput)
    if err != nil {
        return "", fmt.Errorf("failed to get photo: %w", err)
    }

    if getItemOutput.Item == nil {
        return "", ErrPhotoDoesNotExist
    }

    urlAttr, urlOk := getItemOutput.Item["url"].(*types.AttributeValueMemberS)
    ownerAttr, ownerOk := getItemOutput.Item["owner"].(*types.AttributeValueMemberS)

    if !urlOk || !ownerOk {
        return "", fmt.Errorf("missing required attributes in photo item")
    }

    url := urlAttr.Value
    owner := ownerAttr.Value

    if owner != req_id {
        return "", ErrDeletePhotoForbidden
    }

    deleteItemInput := &dynamodb.DeleteItemInput{
        TableName: aws.String("Photo"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: photo_id},
        },
    }

    _, err = db.c.DeleteItem(context.TODO(), deleteItemInput)
    if err != nil {
        return "", fmt.Errorf("failed to delete photo: %w", err)
    }

    return url, nil
}
