/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserDoesNotExist = errors.New("user does not exist")
var ErrBanned = errors.New("user banned")
var ErrPhotoDoesNotExist = errors.New("photo does not exist")
var ErrDeletePhotoForbidden = errors.New("user cannot delete other users' photos")
var ErrNotSelfLike = errors.New("user cannot like his own photos")
var ErrCommentDoesNotExist = errors.New("comment does not exist")
var ErrForbidden = errors.New("operation not allowed")
var ErrUsernameAlreadyTaken = errors.New("username already taken")

// Fountain struct represent a fountain in every API call between this package and the outside world.
// Note that the internal representation of fountain in the database might be different.
type Fountain struct {
	ID        uint64
	Latitude  float64
	Longitude float64
	Status    string
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	DoLogin(username string) (string, error)
	GetUserProfile(id string, req_id string) (User, error)
	CheckBan(user_id string, target_user_id string) (bool, error) // don't know if leave it here or not
	CheckUser(user_id string) (bool, error)
	CheckPhoto(photo_id string) (bool, error)
	FollowUser(user_id string, target_user_id string) error
	BanUser(user_id string, target_user_id string) error
	UnfollowUser(user_id string, target_user_id string) error
	UnbanUser(user_id string, target_user_id string) error
	UploadPhoto(id string, created_at string, url string, owner string) (Photo, error)
	DeletePhoto(photo_id string, req_id string) (string, error)
	LikePhoto(photo_id string, user_id string) error
	UnlikePhoto(photo_id string, user_id string) error
	CommentPhoto(cid string, pid string, uid string, text string, created_datetime string) (Comment, error)
	UncommentPhoto(photo_id string, comment_id string, req_id string) error
	GetMyStream(id string) ([]Photo, error)
	GetUsers(req_id string) ([]UserShortInfo, error)
	SearchUser(username string, req_id string) (UserShortInfo, error)
	GetPhoto(id string, req_id string) (Photo, error)
	SetMyUsername(id string, username string) error
	GetFollowers(id string, req_id string) ([]UserShortInfo, error)
	GetFollowing(id string, req_id string) ([]UserShortInfo, error)
	GetBanned(id string, req_id string) ([]UserShortInfo, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Creating the database tables
	var tableName string
	// Check if table User exists. If not, we need to create it
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='User';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE User(id VARCHAR(32) NOT NULL PRIMARY KEY CHECK (length(id) >= 8), username VARCHAR(16) NOT NULL UNIQUE CHECK (length(username) >= 3));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating User table: %w", err)
		}
	}
	// Check if table Photo exists. If not, we need to create it
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Photo';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Photo(id VARCHAR(32) NOT NULL PRIMARY KEY CHECK (length(id) >= 8), created_at TEXT NOT NULL CHECK (created_at IS datetime(created_at, '+0 days')), url VARCHAR(200) NOT NULL, owner VARCHAR(32), FOREIGN KEY(owner) REFERENCES User(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating Photo table: %w", err)
		}
	}
	// Check if table Like exists. If not, we need to create it
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Like';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Like(photo VARCHAR(32) NOT NULL, user VARCHAR(32) NOT NULL, PRIMARY KEY(photo, user), FOREIGN KEY(photo) REFERENCES Photo(id) ON DELETE CASCADE, FOREIGN KEY(user) REFERENCES User(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating Like table: %w", err)
		}
	}
	// Check if table Comment exists. If not, we need to create it
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Comment';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Comment(id VARCHAR(32) NOT NULL PRIMARY KEY CHECK (length(id) >= 8), photo VARCHAR(32), user VARCHAR(32), text VARCHAR(300) NOT NULL CHECK (length(text) >= 3), created_at TEXT NOT NULL CHECK (created_at IS datetime(created_at, '+0 days')), FOREIGN KEY(photo) REFERENCES Photo(id) ON DELETE CASCADE, FOREIGN KEY(user) REFERENCES User(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating Comment: %w", err)
		}
	}
	// Check if table Following exists. If not, we need to create it
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Following';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Following(user_following VARCHAR(32) NOT NULL, user_followed VARCHAR(32) NOT NULL, PRIMARY KEY(user_following, user_followed), FOREIGN KEY(user_following) REFERENCES User(id), FOREIGN KEY(user_followed) REFERENCES User(id), CHECK (user_following <> user_followed));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating Following: %w", err)
		}
	}
	// Check if table Banned exists. If not, we need to create it
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Banned';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Banned(user_banning VARCHAR(32) NOT NULL, user_banned VARCHAR(32) NOT NULL, PRIMARY KEY(user_banning, user_banned), FOREIGN KEY(user_banning) REFERENCES User(id), FOREIGN KEY(user_banned) REFERENCES User(id), CHECK (user_banning <> user_banned));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure when creating Banned table: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
