package red_test

import (
	"context"
	red "github.com/goclub/redis"
	radix4 "github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type Radix4Client struct {
	radix4.Client
}
func (c Radix4Client) RedisDo(ctx context.Context, valuePtr interface{}, args []string) (result red.Result, err error){
	data := radix4.Maybe{Rcv: valuePtr}
	var moreArg []string
	if len(args) >1 { moreArg = args[1:] }
	err = c.Do(ctx, radix4.Cmd(&data, args[0], moreArg...)) ; if err != nil {
		return
	}
	result = red.Result{
		IsNil: data.Null,
		IsEmpty: data.Empty,
	}
	return
}

func (c Radix4Client)  RedisScript (ctx context.Context, script red.RedisScript) (result red.Result, err error){
	data := radix4.Maybe{Rcv: script.ValuePtr}
	err = c.Do(ctx, radix4.NewEvalScript(script.Script).Cmd(&data, script.Keys, script.Args...)) ; if err != nil {
		return
	}
	result = red.Result{
		IsNil: data.Null,
		IsEmpty: data.Empty,
	}
	return
}
type Test struct {
	T *testing.T
	Expected string
}
func (test Test) RedisDo(ctx context.Context, valuePtr interface{}, args []string) (result red.Result, err error) {
	assert.Equal(test.T, test.Expected, strings.Join(args, " "))
	return
}
func (test Test)  RedisScript (ctx context.Context, script red.RedisScript) (result red.Result, err error){
	return
}