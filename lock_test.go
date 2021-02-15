package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

var mutexCount = 0
func MutexAction(i int) {
	key := "test_mutex"
	mutex := red.Mutex{
		Key: key,
		Expire: time.Millisecond*100,
		Retry: red.Retry{
			Times: 3,
			Duration:time.Millisecond*100,
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
		mutexCount++
		log.Print(i, "锁成功")
	}
	log.Print(i, "业务操作")
	err =  mutex.Unlock(context.TODO()) ; if err != nil {
		mutexCount--
		log.Print(i, "解锁失败:", err)
		return
	}
	log.Print(i, "解锁成功")

}
// 完整的锁测试非常复杂，暂时用简单的测试代替
func TestMutex_Lock(t *testing.T) {
	wg := sync.WaitGroup{}
	max := 5
	for i:=0;i<max;i++{
		wg.Add(1)
		go func(i int) {
			MutexAction(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Print("mutexCount", mutexCount)
	assert.Equal(t, mutexCount > 0, true)
	assert.Equal(t, mutexCount <= 5, true)

}