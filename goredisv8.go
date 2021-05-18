package red

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

func ArgsToInterfaces(args []string) (interfaces []interface{}) {
	for _, s := range args {
		interfaces = append(interfaces, s)
	}
	return
}
type GoRedisV8 struct {
	Core redis.UniversalClient
}
func (r GoRedisV8) DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error) {
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", true, nil
		}
		return
	}
	reply = cmd.Val().(string)
	return
}

func (r GoRedisV8) DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error) {
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, true, nil
		}
		return
	}
	reply = cmd.Val().(int64)
	return
}

func (r GoRedisV8) DoIntegerSliceReply(ctx context.Context, args []string)(reply []int64, isNil bool, err error) {
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, true, nil
		}
		return
	}
	values := cmd.Val().([]interface{})
	for _, v := range values {
		i := v.(int64)
		reply = append(reply, i)
	}
	return
}

func (r GoRedisV8) Eval(ctx context.Context, data Script) (reply interface{}, isNil bool, err error) {
	var argv []interface{}
	for _, s := range data.Argv {
		argv = append(argv, s)
	}
	cmd := r.Core.Eval(ctx, data.Script, data.Keys, argv...)
	err = cmd.Err() ; if err != nil {
		if errors.Is(err, redis.Nil) {
			return "",true, nil
		} else {
			return "",false, err
		}
	}
	reply = cmd.Val()
	return
}