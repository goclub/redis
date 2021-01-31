package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"log"
	"sync"
	"testing"
	"time"
)

func MutexAction(i int) {
	key := "test_mutex"
	mutex := red.Mutex{
		Key: key,
		Expires: time.Second*10,
		Retry: red.Retry{
			Times: 10,
			Duration:time.Second,
		},
	}
	lockSuccess, err := mutex.Lock(context.TODO(), radixClient) ; if err != nil {
		log.Print(i, "锁失败")
		return
	}
	if lockSuccess == false {
		log.Print(i, "锁被占用")
		return
	}
	log.Print(i, "锁成功")
	time.Sleep(time.Second*1)
	log.Print(i, "业务操作")
	unlockOk, err :=  mutex.Unlock(context.TODO()) ; if err != nil {
		log.Print(i, err)
		return
	}
	if unlockOk == false {
		log.Print(i, "撤销业务操作")
	}
}
func TestMutex_Lock(t *testing.T) {
	wg := sync.WaitGroup{}
	for i:=0;i<5;i++{
		wg.Add(1)
		go func(i int) {
			MutexAction(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}