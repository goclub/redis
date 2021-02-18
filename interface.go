package red

import (
	"context"
	redScript "github.com/redis-driver/script"
)

type Client interface {
	RedisCommand (ctx context.Context, valuePtr interface{}, args []string) (result struct {IsNil bool}, err error)
	RedisScript (
		ctx context.Context,
		data redScript.Script,
	) (result struct {IsNil bool}, err error)
	Close () error
}
