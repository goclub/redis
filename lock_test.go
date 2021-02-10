package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func MutexAction(i int) {
	key := "test_mutex"
	mutex := red.Mutex{
		Key: key,
		Expire: time.Second*10,
		Retry: red.Retry{
			Times: 10,
			Duration:time.Second/5,
		},
	}
	lockSuccess, err := mutex.Lock(context.TODO(), radixClient) ; if err != nil {
		log.Print(i, "锁失败")
		return
	}
	if lockSuccess == false {
		log.Print(i, "锁被占用")
		return
	} else {
		log.Print(i, "锁成功")
	}
	time.Sleep(time.Second)
	_, err = red.RPUSH{Key: "test_mutex_list", Value: strconv.Itoa(i)}.Do(context.TODO(), radixClient) ; if err != nil {
		log.Print(err)
		return
	}
	log.Print(i, "业务操作")
	unlockOk, err :=  mutex.Unlock(context.TODO()) ; if err != nil {
		log.Print(i, err)
		return
	}
	if unlockOk == false {
		log.Print(i, "撤销业务操作")
	} else {
		log.Print(i, "解锁成功")
	}
}
func TestMutex_Lock(t *testing.T) {
	_, err := red.DEL{Key:"test_mutex_list"}.Do(context.TODO(), radixClient) ; if err != nil {
		log.Print(err)
		return
	}
	wg := sync.WaitGroup{}
	for i:=0;i<5;i++{
		wg.Add(1)
		go func(i int) {
			MutexAction(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	list, err := red.LRANGE{Key: "test_mutex_list", Start: 0, Stop: -1}.Do(context.TODO(), radixClient) ; if err != nil {
		log.Print(err)
		return
	}
	log.Print("list", list)
}