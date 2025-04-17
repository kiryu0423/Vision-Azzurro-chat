package notify

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func SubscribeToUser(rdb *redis.Client, userID uint, sendToClient func(payload string)) {
	channel := "user:" + strconv.Itoa(int(userID))
	pubsub := rdb.Subscribe(context.Background(), channel)

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			log.Println("[通知受信]", msg.Payload)
			sendToClient(msg.Payload) // WebSocketなどに送信
		}
	}()
}
