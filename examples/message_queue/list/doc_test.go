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
	client, err := NewClient() ; if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	for i:=0;i<10;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := production(client, xtest.UUID(), xtest.Name()) ; if err != nil {
				log.Print(err)
			}
		}()
	}
	go func() {
		consume(client)
	}()
	wg.Wait()
	select{}
}

func production(client red.Client, id string, name string) (err error) {
	log.Print("production", id, name)
	_, err = red.RPUSH{Key: "sendWelcomeMessagePending", Value: id + ":" + name}.Do(context.TODO(), client) ; if err != nil {
		return
	}
	return nil
}

func consume (client red.Client) {
	for {
		payload, isNil, err := red.BRPOPLPUSH{
			Source: "sendWelcomeMessagePending",
			Destination: "sendWelcomeMessageProcessing",
		}.Do(context.TODO(), client) ; if err != nil {
			// 错误不panic 可能是超长时间链接中断错误
			log.Print(err)
		}
		log.Print("consume", payload, isNil)
	}
}

