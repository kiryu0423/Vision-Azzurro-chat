package notify

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func PublishToUser(rdb *redis.Client, userID uint, data map[string]interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	channel := "user:" + strconv.Itoa(int(userID))
	return rdb.Publish(ctx, channel, payload).Err()
}
