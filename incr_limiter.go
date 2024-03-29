package red

import (
	"context"
	xerr "github.com/goclub/error"
	"strconv"
	"time"
)

// IncrLimiter 递增限制器
// eg:消息队列:重新入队:消息ID作为key在1分钟内只能递增3次.三次内返回的Limited为false，超过三次为true
type IncrLimiter struct {
	Key    string        	`note:"key" eg:"mq_requeue:{messageID}"`
	Expire time.Duration 	`note:"有效期" eg:"time.Minute"`
	Maximum   uint64        `note:"最大限制" eg:"3"`
	Increment uint64        `note:"递增值" default:"1"`
}

func (v IncrLimiter) Do(ctx context.Context, client Connecter) (limited bool, err error) {
	// 参数校验
	if v.Key == "" {
		return false, xerr.New("goclub/redis: IncrLimiter{}.Key can not be empty string")
	}
	if v.Expire.Milliseconds() < 1 {
		return false, xerr.New("goclub/redis: IncrLimiter{}.Expire can not less 1 millisecond")
	}
	if v.Maximum < 1 {
		return false, xerr.New("goclub/redis: IncrLimiter{}.Maximum can not less 1")
	}
	if v.Increment == 0 { v.Increment = 1 }
	// 递增脚本
	var isNil bool
	_, isNil, err = client.Eval(ctx, Script{
		KEYS: []string{
			/*1*/ v.Key,
		},
		ARGV: []string{
			/*1*/ strconv.FormatInt(v.Expire.Milliseconds(), 10),
			/*2*/ strconv.FormatUint(v.Maximum, 10),
			/*3*/ strconv.FormatUint(v.Increment, 10),
		},
		Script: `
			local namespace = KEYS[1]
			local expire = ARGV[1]
			local maximun = tonumber(ARGV[2])
			local increment =  tonumber(ARGV[3])
			
			local num = redis.call("GET", namespace)
			if num ~= false and tonumber(num) + increment > maximun then
				return false
			end

			local newNum = redis.call("INCRBY", namespace, increment)
			if newNum == increment then
				redis.call("PEXPIRE", namespace, expire)
			end
			return "OK"
		`,
	})
	// 没有成功递增
	if isNil {
		limited = true
		return
	}
	return
}
