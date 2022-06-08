package red

import (
	"github.com/go-redis/redis/v8"
)

var Connecters = []Connecter{}

func init () {
	Connecters = append(Connecters, GoRedisV8{
		Core: redis.NewClient(&redis.Options{
			// 设置为 15 是为了万一在正式环境运行了测试导致破坏了数据
			DB: 15,
		}),
	})

}
