package red

import "context"

type Doer interface {
	RedisDo (ctx context.Context, valuePtr interface{}, args []string) (result Result, err error)
}