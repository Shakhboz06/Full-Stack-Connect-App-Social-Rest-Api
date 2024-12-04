package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("resource already exists")	
	QueryTimeoutDuration = time.Second * 5

)
type Storage struct {
	Posts interface {
		Create(context.Context, *Posts)  error
		GetbyID(context.Context, int64) (*Posts, error)
		DeletebyID(context.Context, int64) error
		UpdatebyID(context.Context, *Posts) error
		GetUserFeed(context.Context, int64, PaginatedFeed)([]PostswithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *Users) error
		GetUser(context.Context, int64) (*Users, error)
		CreateAndInvite(context.Context, *Users, string, time.Duration)error
		Activation(context.Context, string) error
		Delete(context.Context, int64)error
		GetByEmail(context.Context, string)(*Users, error)
	}
	Comments interface{
		GetbyPostID(context.Context, int64)([]Comment, error)
		Create(context.Context, *Comment)error
	}
	Followers interface{
		Follow(ctx context.Context, FollowerID, userID int64) error
		Unfollow(ctx context.Context, FollowerID, userID int64) error
	}
	Roles interface{
		GetByName(context.Context, string)(*Roles, error)
	}

}

func NewPostgresStorage (db *sql.DB) Storage{
	return Storage{
		Posts: &PostsStore{db},
		Users: &UserStore{db},
		Comments: &CommentsStore{db},
		Followers: &FollowersStore{db},
		Roles: &RolesStore{db},
	}
}



func withTx(db *sql.DB, ctx context.Context, fn func (*sql.Tx) error)error{
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil{
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}