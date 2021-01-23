package red

import "context"

type Doer interface {
	RedisDo (ctx context.Context, valuePtr interface{}, args []string) (result Result, err error)
	RedisScript (ctx context.Context, data RedisScript) (result Result, err error)
}
type RedisScript struct {
	ValuePtr interface{}
	Script string
	Keys []string
	Args []string
}
