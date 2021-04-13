package red

import (
	"github.com/go-redis/redis/v8"
)

var Connecters = []red.Connecter{}

func init () {
	Connecters = append(Connecters, red.GoRedisV8{
		Core: redis.NewClient(&redis.Options{}),
	})
}