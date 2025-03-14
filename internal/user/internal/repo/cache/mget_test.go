package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisMget(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})

	for i := range 10 {
		user := domain.User{
			Id:       int64(i + 1),
			Username: fmt.Sprintf("user%d", i+1),
			Account:  fmt.Sprintf("account%d", i+1),
		}
		val, err := json.Marshal(user)
		require.NoError(t, err)
		err = cmd.Set(context.Background(), fmt.Sprintf("user:%d", user.Id), val, time.Minute).Err()
		require.NoError(t, err)
	}

	res, err := cmd.MGet(context.Background(), "user:1", "user:4", "user:6").Result()
	require.NoError(t, err)
	for _, v := range res {
		if v == "" {
			continue
		}
		var user domain.User
		err = json.Unmarshal([]byte(v.(string)), &user)
		require.NoError(t, err)
		fmt.Println(user)
	}
}
