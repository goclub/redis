package red

import (
	"context"
	"strconv"
	"time"
)

type BRPOPLPUSH struct {
	Source string
	Destination string
	Timeout time.Duration
}
type ResultBRPOPLPUSH struct {
	WaitTime time.Duration
	Element string
}
// 假如在指定时间内没有任何元素被弹出，则返回一个 nil 和等待时长。
// 反之，返回一个含有两个元素的列表，第一个元素是被弹出元素的值，第二个元素是等待时长。
func (data BRPOPLPUSH) Do(ctx context.Context, doer Doer) (value string, isNil bool, err error) {
	cmd := "BRPOPLPUSH"
	err = checkKey(cmd, "Source", data.Source) ; if err != nil {
		return
	}
	err = checkKey(cmd, "Destination", data.Destination) ; if err != nil {
		return
	}
	err = checkDuration(cmd, "Timeout", data.Timeout) ; if err != nil {
		return
	}
	timeoutStr := strconv.FormatInt(int64(data.Timeout/time.Second), 10)
	doResult, err := doer.RedisDo(ctx, &value, []string{cmd, data.Source, data.Destination, timeoutStr,}) ; if err != nil {
		return
	}
	if doResult.IsNil {
		return "", true, nil
	}
	return
}
