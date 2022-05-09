package example_test

import (
	"context"
	red "github.com/goclub/redis"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestEval(t *testing.T) {
	err := func() (err error) {
		ctx := context.Background()
		client, err := NewClient(ctx)
		if err != nil {
			return
		}
		log.Print("使用redis lua 脚本进行计数器递增,一分钟内最多递增3次")
		reply, isNil, err := client.Eval(ctx, red.Script{
			KEYS: []string{
				/* 1 */ "example_count",
			},
			ARGV: []string{
				/* 1 */ "3",
				/* 2 */ strconv.FormatInt(time.Minute.Milliseconds(), 10),
			},
			Script: `
local countKey = KEYS[1]
local threshold = tonumber(ARGV[1])
local expireMill = ARGV[2]
local oldCount = redis.call("GET", countKey)

-- redis 返回的 nil 在 lua 中会被转换为 false
if oldCount ~= false then
	if tonumber(oldCount) >= threshold then
		return false
	end
end
local newIncrValue = tonumber(redis.call("INCR", countKey))
if newIncrValue == 1 then
	redis.call("PEXPIRE", countKey, expireMill)
end
return newIncrValue
`,
		})
		if err != nil {
			return
		}
		// lua 返回的 false 在 redis 中会被转换为 nil
		if isNil {
			log.Print("递增失败计数达到最大限制")
		} else {
			log.Print("递增成功,新值是:")
			log.Print(reply.Int64())
		}

		return
	}()
	if err != nil {
		log.Printf("%+v", err)
	}
}
