package red

import "context"

func Command(ctx context.Context, doer Doer, valuePtr interface{}, args []string) (result Result, err error) {
	return doer.RedisCommand(ctx, valuePtr, args)
}