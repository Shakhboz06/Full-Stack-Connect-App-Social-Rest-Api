package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-project/internal/store"
	"time"

	"github.com/go-redis/redis/v8"
)

const UserExpTime = time.Minute

type UsersStore struct {
	rdb *redis.Client
}

func(s *UsersStore) Get(ctx context.Context, userID int64)(*store.Users, error){
	cacheKey := fmt.Sprintf("user-%d", userID)

	
	data, err := s.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil{
		return nil, nil
	}else if err != nil {
		return nil, err
	}

	var user = &store.Users{}
	if data != ""{
		err := json.Unmarshal([]byte(data), user)
		if err != nil{
			return nil, err
		}
	}

	return user, nil
}

func(s *UsersStore) Set(ctx context.Context, users *store.Users)error{
	if users.ID == 0 {
		return nil
	}

	cacheKey := fmt.Sprintf("user-%v", users.ID)

	json, err := json.Marshal(users)

	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, json, UserExpTime).Err()
}
