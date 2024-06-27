package database_nosql

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    
    "errors"
	"fmt"
	"context"
)

var ErrUserDoesNotExist = errors.New("user does not exist")
var ErrCannotUnfollow = errors.New("user not found in following list")
var ErrCannotUnban = errors.New("user not found in banned list")
var ErrBanned = errors.New("user banned")
var ErrPhotoDoesNotExist = errors.New("photo does not exist")
var ErrDeletePhotoForbidden = errors.New("user cannot delete other users' photos")
var ErrNotSelfLike = errors.New("user cannot like his own photos")
var ErrCommentDoesNotExist = errors.New("comment does not exist")
var ErrForbidden = errors.New("operation not allowed")
var ErrUsernameAlreadyTaken = errors.New("username already taken")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
    DoLogin(username string) (string, error)
    CheckBan(user_id string, target_user_id string) (bool, error)
    CheckUser(user_id string) (bool, error)
    CheckPhoto(photo_id string) (bool, error)
    FollowUser(user_id string, target_user_id string) error
    UnfollowUser(user_id string, target_user_id string) error
    BanUser(user_id string, target_user_id string) error
    UnbanUser(user_id string, target_user_id string) error
    GetFollowing(id string, req_id string) ([]UserShortInfo, error)
    GetBanned(id string, req_id string) ([]UserShortInfo, error)
    GetFollowers(id string, req_id string) ([]UserShortInfo, error)
    SetMyUsername(id string, username string) error
    SearchUser(username string, req_id string) (UserShortInfo, error)
    GetUsers(req_id string) ([]UserShortInfo, error)
    UploadPhoto(id string, created_at string, url string, owner string) (Photo, error)
    DeletePhoto(photo_id string, req_id string) (string, error)
    /*
    GetUserProfile(id string, req_id string) (User, error)
    LikePhoto(photo_id string, user_id string) error
    UnlikePhoto(photo_id string, user_id string) error
    CommentPhoto(cid string, pid string, uid string, text string, created_datetime string) (Comment, error)
    UncommentPhoto(photo_id string, comment_id string, req_id string) error
    GetMyStream(id string) ([]Photo, error)
    GetPhoto(id string, req_id string) (Photo, error)
    */

    // Ping checks whether the database is available or not (in that case, an error will be returned)
    Ping() error
}

type appdbimpl struct {
    c *dynamodb.Client
}

/*
// getItem returns an item if found based on the key provided.
// the key could be either a primary or composite key and values map.
func GetItem(c *dynamodb.Client, tableName string, key DynoNotation) (item DynoNotation, err error) {
    resp, err := c.GetItem(context.TODO(), &dynamodb.GetItemInput{Key: key, TableName: aws.String(tableName)})
    if err != nil {
        return nil, err
    }

    return resp.Item, nil
}
*/

func New(db *dynamodb.Client) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	ctx := context.TODO()

    // Helper function to create a table if it doesn't exist
    createTable := func(tableName string, attributes []types.AttributeDefinition, keySchema []types.KeySchemaElement, provisionedThroughput *types.ProvisionedThroughput) error {
        _, err := db.DescribeTable(ctx, &dynamodb.DescribeTableInput{
            TableName: aws.String(tableName),
        })
        if err != nil {
            var nfe *types.ResourceNotFoundException
            if errors.As(err, &nfe) {
                _, err := db.CreateTable(ctx, &dynamodb.CreateTableInput{
                    TableName:            aws.String(tableName),
                    AttributeDefinitions: attributes,
                    KeySchema:            keySchema,
                    ProvisionedThroughput: provisionedThroughput,
                })
                if err != nil {
                    return fmt.Errorf("error creating table %s: %w", tableName, err)
                }
            } else {
                return err
            }
        }
        return nil
    }

    // Check and create User table if it doesn't exist
    err := createTable("User",
        []types.AttributeDefinition{
            {AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
        },
        []types.KeySchemaElement{
            {AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
        },
        &types.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
    )
    if err != nil {
        return nil, err
    }

    // Check and create Photo table if it doesn't exist
    err = createTable("Photo",
        []types.AttributeDefinition{
            {AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
        },
        []types.KeySchemaElement{
            {AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
        },
        &types.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
    )
    if err != nil {
        return nil, err
    }

	return &appdbimpl{
        c: db,
    }, nil
}

func (db *appdbimpl) Ping() error {
    // DynamoDB does not have a direct Ping equivalent.
    // We can list tables as a simple way to check connectivity.
    _, err := db.c.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
    return err
}
