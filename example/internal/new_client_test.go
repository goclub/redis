package example_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	red "github.com/goclub/redis"
)

func NewClient (ctx context.Context) (client red.GoRedisV8, err error) {
	client = red.GoRedisV8{
		Core: redis.NewClient(&redis.Options{}),
	}
	_, err = client.DoStringReplyWithoutNil(ctx, []string{"PING"}) ; if err != nil {
	    return
	}
	return
}
