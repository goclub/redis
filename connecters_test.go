package red

import (
	"github.com/go-redis/redis/v8"
)

var Connecters = []Connecter{}

func init () {
	Connecters = append(Connecters, GoRedisV8{
		Core: redis.NewClient(&redis.Options{}),
	})
}