package red

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"time"
)
// 使用 GET 时注意考虑要获取的 key 是否存在并发问题
// 不要用 GET 去获取锁，而是通过 SETNX 获取锁(严谨的 SETNX 锁比较复杂，建议 red.DistLock )
type GET struct {
	Key string
}
func (data GET) Do(ctx context.Context, doer Doer) (value string, hasValue bool ,err error) {
	if len(data.Key) == 0 {
		return "",false, errors.New("goclub/redis: GET{} Key cannot be empty")
	}
		args := []string{"GET", data.Key}
	result, err := doer.RedisCommand(ctx, &value, args) ; if err != nil {
		return "", false, err
	}
	if result.IsNil {
		return "", false, err
	} else {
		return value, true, nil
	}
}
type SET struct {
	Key string
	Value string
	Expires time.Duration
}
func (data SET) Do(ctx context.Context, doer Doer) (err error) {
	_, err = coreSET{
		Key: data.Key,
		Value: data.Value,
		Expires: data.Expires,
	}.Do(ctx, doer) ; if err != nil {
		return
	}
	return
}
type SETNX struct {
	Key string
	Value string
	Expires time.Duration
}
func (data SETNX) Do(ctx context.Context, doer Doer) (ok bool,err error) {
	result, err := coreSET{
		Key: data.Key,
		Value: data.Value,
		Expires: data.Expires,
		NX: true,
	}.Do(ctx, doer) ; if err != nil {
		return
	}
	if result.IsNil {
		return false, nil
	} else {
		return true, nil
	}
}
type SETXX struct {
	Key string
	Value string
	Expires time.Duration
}
func (data SETXX) Do(ctx context.Context, doer Doer) (ok bool,err error) {
	result, err := coreSET{
		Key: data.Key,
		Value: data.Value,
		Expires: data.Expires,
		XX: true,
	}.Do(ctx, doer) ; if err != nil {
		return
	}
	if result.IsNil {
		return false, nil
	} else {
		return true, nil
	}
}
type coreSET struct {
	Key string
	Value string
	Expires time.Duration
	NX bool
	XX bool
}
func (data coreSET) Do(ctx context.Context, doer Doer) (result Result,err error) {
	if len(data.Key) == 0 {
		return result, errors.New("goclub/redis: SET{} Key cannot be empty")
	}
	args := []string{"SET", data.Key, data.Value}
	if data.Expires != 0 {
		px := strconv.FormatInt(int64(data.Expires / time.Millisecond), 10)
		args = append(args, "PX", px)
	}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	return doer.RedisCommand(ctx, nil, args)
}
type DEL struct {
	Key string
	Keys []string
}
func (data DEL) Do(ctx context.Context, doer Doer) (delCount uint, err error) {
	args := []string{"DEL"}
	if len(data.Key) != 0 {
		data.Keys = append(data.Keys, data.Key)
	}
	if len(data.Keys) == 0 {
		return 0, errors.New("goclub/redis: DEL{} Keys cannot be empty")
	}
	args = append(args, data.Keys...)
	_, err = doer.RedisCommand(ctx, &delCount, args) ; if err != nil {
		return
	}
	return
}
// 不要使用 DECR 做库存超卖，因为 DECR -1 后的 INCR 无法保证原子性，
// 在需要使用 INCR 撤销 DECR 时可能因为网络进程等各种原因导致执行失败
// 超卖问题建议使用 Redis Lua 脚本保证原子性来解决
type DECR struct {
	Key string
}
func (data DECR) Do(ctx context.Context, doer Doer) (value int64 ,err error) {
	if len(data.Key) == 0 {
		return 0, errors.New("goclub/redis: DECR{} Key cannot be empty")
	}
	args := []string{"DECR", data.Key}
	_, err = doer.RedisCommand(ctx, &value, args) ; if err != nil {
		return 0, err
	}
	return
}
type INCR struct {
	Key string
}
func (data INCR) Do(ctx context.Context, doer Doer) (value int64 ,err error) {
	if len(data.Key) == 0 {
		return 0, errors.New("goclub/redis: INCR{} Key cannot be empty")
	}
	args := []string{"INCR", data.Key}
	_, err = doer.RedisCommand(ctx, &value, args) ; if err != nil {
		return 0, err
	}
	return
}