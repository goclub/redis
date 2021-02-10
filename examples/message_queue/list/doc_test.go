package doc_test

import (
	"context"
	red "github.com/goclub/redis"
	xtest "github.com/goclub/test"
	"log"
	"sync"
	"testing"
)

func TestMessageQueue(t *testing.T) {
	doer, err := NewClient() ; if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	for i:=0;i<10;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := production(doer, xtest.UUID(), xtest.Name()) ; if err != nil {
				log.Print(err)
			}
		}()
	}
	go func() {
		consume(doer)
	}()
	wg.Wait()
	select{}
}

func production(doer red.Doer, id string, name string) (err error) {
	log.Print("production", id, name)
	_, err = red.RPUSH{Key: "sendWelcomeMessagePending", Value: id + ":" + name}.Do(context.TODO(), doer) ; if err != nil {
		return
	}
	return nil
}

func consume (doer red.Doer) {
	for {
		payload, isNil, err := red.BRPOPLPUSH{
			Source: "sendWelcomeMessagePending",
			Destination: "sendWelcomeMessageProcessing",
		}.Do(context.TODO(), doer) ; if err != nil {
			// 错误不panic 可能是超长时间链接中断错误
			log.Print(err)
		}
		log.Print("consume", payload, isNil)
	}
}

