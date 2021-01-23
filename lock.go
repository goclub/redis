package red

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

type Mutex struct {
	Key string
	Expires time.Duration
	RetryCount uint8
	RetryInterval time.Duration
	startTime time.Time
	lockValue string
	doer Doer
}

func (data Mutex) Unlock (ctx context.Context) (unlockOk bool ,err error) {
	if data.startTime.After(time.Now().Add(data.Expires)) {
		return false, nil
	}
	var delCount uint
	script := `
if redis.call("GET", KEYS[1] == ARGV[1])
then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`
	_, err = data.doer.RedisScript(ctx, RedisScript{
		ValuePtr: &delCount,
		Script: script,
		Keys: []string{data.Key},
		Args: []string{data.lockValue},
	}) ; if err != nil {
		return
	}
	switch delCount {
	case 0:
		return false, errors.New("unlock fail, can not found lock key(" + data.Key + ")")
	case 1:
		return true, nil
	default:
		return false, errors.New("del count:" + strconv.Itoa(int(delCount)))
	}
}
func (data Mutex) Lock(ctx context.Context, doer Doer) ( ok bool, err error) {
	data.startTime = time.Now() // start time 必须在 SETNX 之前记录,否则会在SETNX 延迟时候导致时间错误
	data.doer = doer
	data.lockValue = time.Now().String()
	ok, err = SETNX{
		Key: data.Key,
		Value: data.lockValue,
		Expires: data.Expires,
	}.Do(ctx, doer) ; if err != nil {
		return
	}
	return
}
func a()  {
	mutex := Mutex{
		Key: "some",
		Expires: time.Second*10,
	}
	ok, err := mutex.Lock(context.TODO(), nil) ; if err != nil {
		panic(err)
	}
	if ok == false {
		log.Print("锁被占用")
		return
	}
	// do some
	unlockOk, err := mutex.Unlock(context.TODO()) ; if err != nil {
		panic(err)
	}
	if unlockOk == false {
		// rollback
	}
}