package red

import (
	"context"
	"github.com/go-redis/redis/v8"
	xerr "github.com/goclub/error"
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
func NewGoRedisV8(goredisV8Client redis.UniversalClient) GoRedisV8 {
	return GoRedisV8{
		Core:  goredisV8Client,
	}
}
func (r GoRedisV8) DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error) {
	defer func() {
		if err != nil { err = xerr.WithStack(err) }
	}()
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		if xerr.Is(err, redis.Nil) {
			return "", true, nil
		}
		return
	}
	reply, err = Reply{cmd.Val()}.String() ; if err != nil {
	    return
	}
	return
}
func (r GoRedisV8) DoStringReplyWithoutNil(ctx context.Context, args []string) (reply string, err error) {
	var isNil bool
	reply, isNil, err = r.DoStringReply(ctx, args) ; if err != nil {
	    return
	}
	if isNil == true {
		err = xerr.New("DoStringReplyWithoutNil(ctx, args) args exec result can not be nil")
		return
	}
	return
}

func (r GoRedisV8) DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error) {
	defer func() {
		if err != nil { err = xerr.WithStack(err) }
	}()
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		if xerr.Is(err, redis.Nil) {
			return 0, true, nil
		}
		return
	}
	reply, err = Reply{cmd.Val()}.Int64() ; if err != nil {
		return
	}
	return
}
func (r GoRedisV8) DoIntegerReplyWithoutNil(ctx context.Context, args []string) (reply int64, err error) {
	var isNil bool
	reply, isNil, err = r.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	if isNil == true {
		err = xerr.New("DoIntegerReplyWithoutNil(ctx, args) args exec result can not be nil")
		return
	}
	return
}
func (r GoRedisV8) DoArrayIntegerReply(ctx context.Context, args []string)(reply []OptionInt64, err error) {
	defer func() {
		if err != nil { err = xerr.WithStack(err) }
	}()
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		return
	}
	return Reply{cmd.Val()}.Int64Slice()
}

func (r GoRedisV8) DoArrayStringReply(ctx context.Context, args []string)(reply []OptionString, err error) {
	defer func() {
		if err != nil { err = xerr.WithStack(err) }
	}()
	cmd := r.Core.Do(ctx, ArgsToInterfaces(args)...)
	err = cmd.Err() ; if err != nil {
		return
	}
	return Reply{cmd.Val()}.StringSlice()
}

func (r GoRedisV8) Eval(ctx context.Context, data Script) (reply Reply, isNil bool, err error) {
	defer func() {
		if err != nil { err = xerr.WithStack(err) }
	}()
	var argv []interface{}
	for _, s := range data.ARGV {
		argv = append(argv, s)
	}
	cmd := r.Core.Eval(ctx, data.Script, data.KEYS, argv...)
	err = cmd.Err() ; if err != nil {
		if xerr.Is(err, redis.Nil) {
			return Reply{},true, nil
		} else {
			return Reply{},false, err
		}
	}
	reply = Reply{cmd.Val()}
	return
}
func (r GoRedisV8) EvalWithoutNil(ctx context.Context, data Script) (reply Reply, err error) {
	var isNil bool
	reply, isNil, err = r.Eval(ctx, data) ; if err != nil {
		return
	}
	if isNil == true {
		err = xerr.New("EvalWithoutNil(ctx, data) script exec result can not be nil")
		return
	}
	return
}