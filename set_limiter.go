package red

import (
	"context"
	xerr "github.com/goclub/error"
	"strconv"
	"time"
)
// SetLimiter 集合限制器
// 使用场景: 限制用户每天只能试读3个章节(如果不允许一天内反复试读相同章节则可以使用 IncrLimiter )
// 注意:
// 如果 key = free_trial:{userID} 			  Expire = 24h 是限制24小时
// 如果 key = free_trial:2022-01-01:{userID}   Expire = 24h 是限制每天
type SetLimiter struct {
	Key      string        `eg:"free_trial:2022-01-01:{userID}"`
	Member   string         `eg:"{chapterID}"`
	Expire   time.Duration `note:"有效期" eg:"time.Hour*24"`
	Maximum  uint64        `note:"最大限制" eg:"3"`
}

func (v SetLimiter) Do(ctx context.Context, client Connecter) (limited bool, err error) {
	// 参数校验
	if v.Key == "" {
		return false, xerr.New("goclub/redis: SetLimiter{}.Key can not be empty string")
	}
	if v.Member == "" {
		return false, xerr.New("goclub/redis: SetLimiter{}.Member can not be empty string")
	}
	if v.Expire.Milliseconds() < 1 {
		return false, xerr.New("goclub/redis: SetLimiter{}.Expire can not less 1 millisecond")
	}
	if v.Maximum < 1 {
		return false, xerr.New("goclub/redis: SetLimiter{}.Maximum can not less 1")
	}
	// 递增脚本
	var isNil bool
	_, isNil, err = client.Eval(ctx, Script{
		KEYS: []string{
			/*1*/ v.Key,
		},
		ARGV: []string{
			/*1*/ v.Member,
			/*2*/ strconv.FormatInt(v.Expire.Milliseconds(), 10),
			/*3*/ strconv.FormatUint(v.Maximum, 10),
		},
		Script: `
			local namespace = KEYS[1]
			local member = ARGV[1]
			local expire = ARGV[2]
			local maximun = tonumber(ARGV[3])	
			
			local exist = redis.call("SISMEMBER", namespace, member)
			if exist == 1 then
				return "OK"
			end

			local num = redis.call("SCARD", namespace)
			if num ~= false and tonumber(num) >= maximun then
				return false
			end

			local newNum = redis.call("SADD", namespace, member)
			if num == 0 then
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