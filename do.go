package red

import "context"

func Do(ctx context.Context, doer Doer, valuePtr interface{}, args []string) (result Result, err error) {
	return doer.RedisDo(ctx, valuePtr, args)
}