package red

import (
	"context"
	xerr "github.com/goclub/error"
)

type HDEL struct {
	Key   string
	Field []string
}

func (data HDEL) Do(ctx context.Context, client Connecter) (delTotal uint64, err error) {
	if len(data.Key) == 0 {
		err = xerr.New("goclub/redis: key can not be empty string")
		return
	}
	args := []string{"HDEL", data.Key}

	if len(data.Field) == 0 {
		err = xerr.New("goclub/redis: HDEL fields can not be empty")
		return
	}
	args = append(args, data.Field...)
	var value int64
	value, _, err = client.DoIntegerReply(ctx, args)
	if err != nil {
		return
	}
	delTotal = uint64(value)
	return
}
