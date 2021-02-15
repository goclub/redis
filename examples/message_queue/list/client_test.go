package doc_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
)

type radixClient struct {
	radix.Client
}
func (c radixClient) RedisCommand(ctx context.Context, valuePtr interface{}, args []string) (result red.Result, err error){
	data := radix.Maybe{Rcv: valuePtr}
	var moreArg []string
	if len(args) >1 { moreArg = args[1:] }
	err = c.Do(ctx, radix.Cmd(&data, args[0], moreArg...)) ; if err != nil {
		return
	}
	result = red.Result{
		IsNil: data.Null,
		IsEmpty: data.Empty,
	}
	return
}

func (c radixClient)  RedisScript (ctx context.Context, script red.RedisScript) (result red.Result, err error){
	data := radix.Maybe{Rcv: script.ValuePtr}
	err = c.Do(ctx, radix.NewEvalScript(script.Script).Cmd(&data, script.Keys, script.Args...)) ; if err != nil {
		return
	}
	result = red.Result{
		IsNil: data.Null,
		IsEmpty: data.Empty,
	}
	return
}

func NewClient() (red.Client, error) {
	coreClient, err := (radix.PoolConfig{}).New(context.TODO(), "tcp", "127.0.0.1:6379") ; if err != nil {
		return nil, err
	}
	return radixClient{coreClient}, nil
}

