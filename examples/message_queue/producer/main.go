package main

import (
	"context"
	red "github.com/goclub/redis"
	exampleMQData "github.com/goclub/redis/examples/message_queue/data"
	xtest "github.com/goclub/test"
	"github.com/mediocregopher/radix/v4"
	radix4 "github.com/redis-driver/mediocregopher-radix-v4"
	"log"
)

func main() {
	var client radix4.Client
	ctx := context.Background()
	{
		core, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
			log.Print(err) ; return
		}
		client.Core = core
	}
	defer client.Close()
	message := exampleMQData.UserSignIn{
		UserID: xtest.UUID(),
	}
	data, err := red.StructFieldValues(message) ; if err != nil {
		log.Print(err) ; return
	}
	streamID, err := red.XADD{
		Key: message.StreamKey(),
		FieldValues: data,
	}.Do(ctx, client) ; if err != nil {
		log.Print(err) ; return
	}
	log.Print("StreamID:", streamID)
}
