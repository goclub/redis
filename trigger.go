package red

import (
	"context"
	"fmt"
	xerr "github.com/goclub/error"
	"strconv"
	"time"
)

// Trigger 触发器
// 每5分钟出现3次则触发,但是10分钟内只触发一次
// func exmaple () {
// 	triggered, err := Trigger{
// 		Namespace: "pay_fail_alarm",
// 		Interval: time.Minute*5,
// 		Threshold: 3,
// 		Frequency: time.Minute*5,
// 	}.Do(ctx, client) ; if err != nil {
// 	    return
// 	}
// 	if triggered {
// 		// do some
// 	}
// }
type Trigger struct {
	Namespace string         `note:"命名空间" eg:"alarm_login_fail:user:1"`
	Interval  time.Duration  `note:"持续多久"`
	Threshold uint64         `note:"累计多少次"`
	Frequency time.Duration  `note:"多久最多只触发一次 建议与Interval相同"`
}
func (v Trigger) Do(ctx context.Context, client Connecter) (triggered bool, err error) {
	if v.Namespace == "" {
		return false, xerr.New("goclub/redis: Trigger{}.Namespace can not be empty string")
	}
	if v.Interval.Seconds() < 1 {
		return false, xerr.New("goclub/redis: Trigger{}.Interval can not less 1 second")
	}
	if v.Threshold == 0 {
		return false, xerr.New("goclub/redis: Trigger{}.Threshold can not be 0")
	}
	if v.Frequency.Seconds() < 1{
		return false, xerr.New("goclub/redis: Trigger{}.Frequency can not less 1 second")
	}
	reply, err := client.EvalWithoutNil(ctx, Script{
		KEYS:   []string{
			/* 1 */ v.Namespace,
			/* 2 */ v.Namespace + ":frequency",
		},
		ARGV:   []string{
			/* 1 */ strconv.FormatInt(v.Interval.Milliseconds(), 10),
			/* 2 */ strconv.FormatUint(v.Threshold, 10),
			/* 3 */ strconv.FormatInt(v.Frequency.Milliseconds(), 10),
		},
		Script: `
local key       = KEYS[1]
local frequencyKey = KEYS[2]

local intervalMill  = ARGV[1]
local threshold = tonumber(ARGV[2])
local frequencyMill = ARGV[3]

local newCount = tonumber(redis.call("INCR", key))
if newCount == 1 then
	redis.call("PEXPIRE", key, intervalMill)
end
if newCount >= threshold then
	redis.call("DEL", key)
	if redis.call("SET", frequencyKey, 1, "PX", frequencyMill, "NX") == false then
		return 0
	end
	return 1
end
return 0
		`,
	}) ; if err != nil {
	    return
	}
	replyInt, err := reply.Int64() ; if err != nil {
	    return
	}
	switch replyInt {
	case 1:
		triggered = true
	case 0:
		triggered = false
	default:
		err = xerr.New(fmt.Sprintf("goclub/redis: Trigger{}.Do() redis eval reply unexpected", replyInt))
	}
	return
}
