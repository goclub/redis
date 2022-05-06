package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func MutexAction(i int, client Connecter, mutexCount *int,debug bool) {
	key := "test_mutex"
	debugLog := func(v ...interface{}) {
		if debug {
			log.Print(v...)
		}
	}
	mutex := Mutex{
		Key: key,
		Expire: time.Millisecond*100,
		Retry: Retry{
			Times:    3,
			Interval: time.Millisecond*100,
		},
	}
	lockSuccess, unlock, err := mutex.Lock(context.TODO(), client) ; if err != nil {
		debugLog(i, "锁失败")
		return
	}
	if lockSuccess == false {
		debugLog(i, "锁被占用")
		return
	} else {
		*mutexCount++
		debugLog(i, "锁成功")
	}
	debugLog(i, "业务操作")
	err = unlock(context.TODO()) ; if err != nil {
		*mutexCount--
		debugLog(i, "解锁失败:", err)
		return
	}
	debugLog(i, "解锁成功")
}
// 完整的锁测试非常复杂，暂时用简单的测试代替
func TestMutex_Lock(t *testing.T) {
	debug := false
	for _, connecter := range Connecters {
		wg := sync.WaitGroup{}
		var mutexCount = 0
		max := 5
		for i:=0;i<max;i++{
			wg.Add(1)
			go func(i int) {
				MutexAction(i, connecter, &mutexCount, debug)
				wg.Done()
			}(i)
		}
		wg.Wait()
		if debug {
			log.Print("mutexCount", mutexCount)
		}
		assert.Equal(t, mutexCount > 0, true)
		assert.Equal(t, mutexCount <= 5, true)
	}
}