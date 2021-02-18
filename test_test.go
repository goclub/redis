package red_test

import (
	"context"
	redScript "github.com/redis-driver/script"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)
type Test struct {
	T *testing.T
	Expected string
}

func (test Test) RedisCommand(ctx context.Context, valuePtr interface{}, args []string) (result struct {IsNil bool}, err error) {
	assert.Equal(test.T, test.Expected, strings.Join(args, " "))
	return
}
func (test Test) RedisScript (ctx context.Context, script redScript.Script) (result struct {IsNil bool}, err error){
	return
}
func (test Test)  Close () error{
	return nil
}