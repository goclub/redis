package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)
type Test struct {
	T *testing.T
	Expected string
}
func (test Test) RedisCommand(ctx context.Context, valuePtr interface{}, args []string) (result red.Result, err error) {
	assert.Equal(test.T, test.Expected, strings.Join(args, " "))
	return
}
func (test Test)  RedisScript (ctx context.Context, script red.RedisScript) (result red.Result, err error){
	return
}