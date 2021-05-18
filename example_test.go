package red_test

import (
	red "github.com/goclub/redis"
	redis "github.com/go-redis/redis/v8"
)

var exampleClient = red.GoRedisV8{
	Core: redis.NewClient(&redis.Options{}),
}

