package red

import (
	"context"
	xtime "github.com/goclub/time"
	"github.com/pkg/errors"
	"strconv"
	"time"
)
// 使用 GET 时注意考虑要获取的 key 是否存在并发问题
// 不要用 GET 去获取锁，而是通过 SETNX 获取锁(严谨的 SETNX 锁比较复杂，建议 red.DistLock )
type GET struct {
	Key string
}
func (data GET) Do(ctx context.Context, client Client) (value string, hasValue bool ,err error) {
	if len(data.Key) == 0 {
		return "",false, errors.New("goclub/redis:  GET{} Key cannot be empty")
	}
		args := []string{"GET", data.Key}
	result, err := Command(ctx, client, &value, args) ; if err != nil {
		return "", false, err
	}
	if result.IsNil {
		return "", false, err
	} else {
		return value, true, nil
	}
}
// >= 2.6.12: Added the EX, PX, NX and XX options.
type SET struct {
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time
	// >= 6.0: Added the KEEPTTL option.
	KeepTTL bool
	NeverExpire bool
}
func (data SET) Do(ctx context.Context, client Client) (err error) {
	_, err = coreSET{
		Key: data.Key,
		Value: data.Value,
		Expire: data.Expire,
		ExpireAt: data.ExpireAt,
		NeverExpire: data.NeverExpire,
		KeepTTL: data.KeepTTL,
	}.Do(ctx, client) ; if err != nil {
		return
	}
	return
}
type SETNX struct {
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time
	// >= 6.0: Added the KEEPTTL option.
	KeepTTL bool
	NeverExpire bool
}
func (data SETNX) Do(ctx context.Context, client Client) (ok bool,err error) {
	result, err := coreSET{
		Key: data.Key,
		Value: data.Value,
		Expire: data.Expire,
		ExpireAt: data.ExpireAt,
		// >= 6.0: Added the KEEPTTL option.
		KeepTTL: data.KeepTTL,
		NeverExpire: data.NeverExpire,
		NX: true,
	}.Do(ctx, client) ; if err != nil {
		return
	}
	if result.IsNil {
		return false, nil
	} else {
		return true, nil
	}
}
// >= 2.6.12: Added the EX, PX, NX and XX options.
type SETXX struct {
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time
	// >= 6.0: Added the KEEPTTL option.
	KeepTTL bool
	NeverExpire bool
}
func (data SETXX) Do(ctx context.Context, client Client) (ok bool,err error) {
	result, err := coreSET{
		Key: data.Key,
		Value: data.Value,
		Expire: data.Expire,
		ExpireAt: data.ExpireAt,
		KeepTTL: data.KeepTTL,
		NeverExpire: data.NeverExpire,
		XX: true,
	}.Do(ctx, client) ; if err != nil {
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
	Expire time.Duration
	ExpireAt time.Time
	// >= 6.0: Added the KEEPTTL option.
	KeepTTL bool
	NeverExpire bool
	NX bool
	XX bool
}
var ErrSetForgetTimeToLive = errors.New("goclub/redis:  red.SET maybe you forget set field Expire or ExpireAt or KeepTTL")
func (data coreSET) Do(ctx context.Context, client Client) (result Result,err error) {
	if len(data.Key) == 0 {
		return result, errors.New("goclub/redis:  SET{} Key cannot be empty")
	}
	args := []string{"SET", data.Key, data.Value}
	// 只有在明确 NeverExpire 时候才允许 Expire 留空
	if data.NeverExpire == false && data.Expire == 0 && data.ExpireAt.IsZero() && data.KeepTTL == false {
		return Result{}, ErrSetForgetTimeToLive
	}
	if data.Expire != 0 {
		px := strconv.FormatInt(data.Expire.Milliseconds(), 10)
		args = append(args, "PX", px)
	}
	if data.ExpireAt.IsZero() == false {
		args = append(args, "PXAT", strconv.FormatInt(xtime.UnixMilli(data.ExpireAt), 10))
	}
	if data.KeepTTL {
		args = append(args, "KEEPTTL")
	}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	return Command(ctx, client, nil, args)
}
type DEL struct {
	Key string
	Keys []string
}
func (data DEL) Do(ctx context.Context, client Client) (delCount uint, err error) {
	args := []string{"DEL"}
	if data.Key != "" {
		data.Keys = append(data.Keys, data.Key)
	}
	if len(data.Keys) == 0 {
		return 0, errors.New("goclub/redis:  DEL{} Keys cannot be empty")
	}
	args = append(args, data.Keys...)
	_, err = Command(ctx, client, &delCount, args) ; if err != nil {
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
func (data DECR) Do(ctx context.Context, client Client) (value int64 ,err error) {
	if len(data.Key) == 0 {
		return 0, errors.New("goclub/redis:  DECR{} Key cannot be empty")
	}
	args := []string{"DECR", data.Key}
	_, err = Command(ctx, client, &value, args) ; if err != nil {
		return 0, err
	}
	return
}
type INCR struct {
	Key string
}
func (data INCR) Do(ctx context.Context, client Client) (value int64 ,err error) {
	if len(data.Key) == 0 {
		return 0, errors.New("goclub/redis:  INCR{} Key cannot be empty")
	}
	args := []string{"INCR", data.Key}
	_, err = Command(ctx, client, &value, args) ; if err != nil {
		return 0, err
	}
	return
}
type APPEND struct {
	Key string
	Value string
}
func (data APPEND) Do(ctx context.Context, client Client) (length uint, err error) {
	cmd := "APPEND"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	if data.Value == "" {
		err = checkKey(cmd, "Value", data.Value) ; if err != nil {
			return
		}
	}
	args := []string{cmd, data.Key, data.Value}
	_, err = Command(ctx, client, &length, args) ; if err != nil {
		return
	}
	return
}

type GETBIT struct {
	Key string
	Offset OptionUint32
}
func (data GETBIT) Do(ctx context.Context, client Client) (bit uint8, err error) {
	cmd := "GETBIT"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	if data.Offset.valid == false {
		return 0, errForgetArgs(cmd, "offset")
	}
	args := []string{cmd, data.Key, strconv.FormatUint(uint64(data.Offset.Unwrap()), 10)}
	_, err = Command(ctx, client, &bit, args) ; if err != nil {
		return
	}
	return
}

type SETBIT struct {
	Key string
	Offset OptionUint32
	Value OptionUint
}
func (data SETBIT) Do(ctx context.Context, client Client) (oldBit uint8, err error) {
	cmd := "SETBIT"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	if data.Offset.valid == false {
		return 0, errForgetArgs(cmd, "offset")
	}
	offset := strconv.FormatUint(uint64(data.Offset.Unwrap()), 10)
	if data.Value.valid == false {
		return 0, errForgetArgs(cmd, "value")
	}
	value  := strconv.FormatUint(uint64(data.Value.Unwrap()), 10)
	if value != "0" && value != "1" {
		return 0, errors.New("goclub/redis: SETBIT value must be 0 or 1, can not be " + value)
	}
	args := []string{cmd, data.Key, offset, value}
	_, err = Command(ctx, client, &oldBit, args) ; if err != nil {
		return
	}
	return
}