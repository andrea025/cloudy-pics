package database_nosql

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "context"
    "fmt"
)

func (db *appdbimpl) DeletePhoto(photo_id string, req_id string) (string, error) {

    /* CODE TO TRANSLATEFOR DYNAMO DB
    var url, owner string
    err := db.c.QueryRow("SELECT url, owner FROM Photo WHERE id=?", photo_id).Scan(&url, &owner)
    if errors.Is(err, sql.ErrNoRows) {
        return "", ErrPhotoDoesNotExist
    } else if err != nil {
        return "", err
    }

    if owner != req_id {
        return "", ErrDeletePhotoForbidden
    }

    sqlStmt := `DELETE FROM Photo WHERE id == ?;`
    _, err = db.c.Exec(sqlStmt, photo_id)
    return url, err
    */

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
