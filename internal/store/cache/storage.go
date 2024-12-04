package cache

import (
	"context"
	"go-project/internal/store"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64)(*store.Users, error)
		Set(context.Context, *store.Users)error
	
	}
}


func NewRedisStorage(rdb *redis.Client) Storage{
	return Storage{
		Users: &UsersStore{rdb: rdb},
	}
}