package red

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Mutex struct {
	Key string
	Expire time.Duration
	Retry Retry
	startTime time.Time
	lockValue string
	client Client
}

func AsErrUnlock(err error) (unlockErr *ErrUnlock, asErrUnlock bool) {
	asErrUnlock = errors.As(err, &unlockErr)
	return
}
type ErrUnlock struct {
	IsTimeout bool
	IsUnexpectedError bool
	IsConnectErr bool
	Err error
}
// 自定义错误的 Error 方法一定要加 (*Errxxx) 原因：https://github.com/goclub/error
func (e *ErrUnlock) Error() string {
	return e.Err.Error()
}
func (e *ErrUnlock) Unwrap() error {
	return e.Err
}
func (data *Mutex) Unlock (ctx context.Context) (err error) {
	if data.startTime.After(time.Now().Add(data.Expire)) {
		return &ErrUnlock{
			IsTimeout: true,
			Err: errors.New("goclub/redis: IsTimeout Mutex{}.Unlock() key:" + data.Key  + " is timeout"),
		}
	}
	var delCount uint
	script := `
if redis.call("get", KEYS[1]) == ARGV[1]
then
	return redis.call("del", KEYS[1])
else
	return 0
end
`
	_, err = data.client.RedisScript(ctx, RedisScript{
		ValuePtr: &delCount,
		Script: script,
		Keys: []string{data.Key},
		Args: []string{data.lockValue},
	}) ; if err != nil {
		return &ErrUnlock{
			IsConnectErr: true,
			Err: err,
		}
	}
	switch delCount {
	case 0:
		return &ErrUnlock{
			IsTimeout: true,
			Err: errors.New("goclub/redis: IsTimeout Mutex{}.Unlock() key:" + data.Key  + " is timeout"),
		}
	case 1:
		return nil
	default:
		return &ErrUnlock{
			IsUnexpectedError: true,
			Err: errors.New("goclub/redis: IsUnexpectedError Mutex{}.Unlock() del " + data.Key + " count:" + strconv.Itoa(int(delCount))),
		}
	}
}
func (data *Mutex) Lock(ctx context.Context, client Client) ( ok bool, err error) {
	err = data.Retry.check() ; if err != nil {
		return
	}
	retryCount := int(data.Retry.Times)
	return mutexLock(ctx, client, data, &retryCount)
}

func mutexLock(ctx context.Context, client Client, data *Mutex, retryCount *int) (ok bool, err error) {
	data.startTime = time.Now() // start time 必须在 SETNX 之前记录,否则会在SETNX 延迟时候导致时间错误
	data.client = client
	data.lockValue = uuid.NewString()
	ok, err = SETNX{
		Key: data.Key,
		Value: data.lockValue,
		Expire: data.Expire,
	}.Do(ctx, client) ; if err != nil {
		return
	}
	if ok == false {
		*retryCount--
		if *retryCount == -1 {
			return
		}
		time.Sleep(data.Retry.Duration)
		return mutexLock(ctx, client, data, retryCount)
	}
	return
}