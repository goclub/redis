package red

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type LPUSH struct {
	Key string
	Value string
	Values []string
}
func (data LPUSH) Do(ctx context.Context, doer Doer) (length uint, err error) {
	cmd := "LPUSH"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	if len(data.Value) != 0 {
		data.Values = append(data.Values, data.Value)
	}
	args := append([]string{cmd, data.Key}, data.Values...)
	_, err = doer.RedisDo(ctx, &length, args) ; if err != nil {
		return
	}
	return
}

type BRPOPLPUSH struct {
	Source string
	Destination string
	Timeout time.Duration
}
type ResultBRPOPLPUSH struct {
	WaitTime time.Duration
	Element string
}
func (data BRPOPLPUSH) Do(ctx context.Context, doer Doer) (value string, isNil bool, err error) {
	cmd := "BRPOPLPUSH"
	err = checkKey(cmd, "Source", data.Source) ; if err != nil {
		return
	}
	err = checkKey(cmd, "Destination", data.Destination) ; if err != nil {
		return
	}
	if data.Timeout != 0 && data.Timeout < time.Second {
		return "", false, errors.New("goclub/redis:(ERR_TIMEOUT) BRPOPLPUSH Timeout can not less at time.Second")
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
