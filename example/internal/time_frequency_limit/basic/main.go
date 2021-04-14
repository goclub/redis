package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	red "github.com/goclub/redis"
	xtime "github.com/goclub/time"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)
/*
可反复执行 main ，观察执行结果。
*/
func main () {
	client := red.GoRedisV8{
		Core: redis.NewClient(&redis.Options{}),
	}
	ctx := context.Background()
	userID := "u1"
	limited, err := Limited(ctx, client, userID, time.Second * 10) ; if err != nil {
		panic(err)
	}
	log.Print("limited ", limited)
}
func Limited (ctx context.Context, client red.Connecter, userID string, limitDuration time.Duration) (limited bool, err error) {
	recordTimeKey := "example_tfl:" + userID
	nowMilli := strconv.FormatInt(xtime.UnixMilli(time.Now()), 10)
	limitDurationMilli := strconv.FormatInt(limitDuration.Milliseconds(), 10)
	expireMilli := strconv.FormatInt(time.Duration(limitDuration*3).Milliseconds(), 10)
	result, isNil , err := red.Script{
		Keys: []string{
			/* 1 */ recordTimeKey,
		},
		Argv: []string{
			/* 1 */ nowMilli,
			/* 2 */ limitDurationMilli,
			/* 3 */ expireMilli,
		},
		Script: `
			-- 定义变量，便于阅读
			local recordTimeKey = KEYS[1]
			local nowMilli = tonumber(ARGV[1])
			local limitDurationMilli = tonumber(ARGV[2])
			local expireMilli = tonumber(ARGV[3])

			local recordTime =  redis.call('GET', recordTimeKey)
			-- nil 需要通过 == false 判断，而不是 == nil
			local noRecord = recordTime == false
			-- 如果没有记录，则保存记录并通过
			if noRecord then
				-- 设置PX过期时间，以节省 redis 内存空间。
				redis.call("SET", recordTimeKey, nowMilli, "PX", limitDurationMilli)
				return "pass"
			end
			-- 当前记录的时间与已存在的时间的间隔大于限制间隔，则保存新记录并返回不限制
			if nowMilli - recordTime > limitDurationMilli then  
				redis.call("SET", recordTimeKey, nowMilli, "PX", limitDurationMilli)
				return "pass"
			end
			-- 其他情况则限制
			return "limited"
		`,
	}.Do(ctx, client) ; if err != nil {
		return
	}
	if isNil {
		// 代码出现一行中只有 return 会导致 nil,当前script 不存在空return
		return true, errors.New("can not be nil")
	}
	resultString := result.(string)
	switch resultString {
	case "pass":
		limited = false
		return
	case "limited":
		limited = true
		return
	default:
		// 防御编程
		return true, errors.New("Unknown results (" +  resultString + ")")
	}
}
