package main

import (
	"context"
	xhttp "github.com/goclub/http"
	examplesMQEmail "github.com/goclub/redis/examples/message_queue/welcome_email/email"
	exmapleMQ "github.com/goclub/redis/examples/message_queue/welcome_email/mesasge_queue"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	radixClient, err := exmapleMQ.ConnectRedis() ; if err != nil {
		panic(err)
	}
	go func() {
		err = examplesMQEmail.SubscribeUserSignInMessage(ctx, radixClient) ; if err != nil {
			panic(err)
		}
	}()
	xhttp.GracefulClose(func() {
		cancelCtx()
		log.Print("job exit")
	})
}
