package examplesMQEmail

import (
	"context"
	red "github.com/goclub/redis"
	exmapleMQ "github.com/goclub/redis/examples/message_queue/welcome_email/mesasge_queue"
	"log"
	"time"
)
func SyncSendEmail(name string) {
	time.Sleep(time.Second) // 模拟网络io
	log.Print("email: welcome " + name)
}

var emailUserSignInGroupName = "email"
func SubscribeUserSignInMessage(ctx context.Context, radixClient red.DriverRadixClient4) (error) {
	key := exmapleMQ.UserSignInMessage{}.StreamKey()
	group := "email"
	err := red.XGROUPCreate{
		Key: key,
		Group: group,
		ID: "0",
		MKSTREAM: true,
	}.Do(ctx, radixClient) ; if err != nil {
		log.Print(err)
	}
	for {
		select {
			case <- ctx.Done():
				return ctx.Err()
		default:
			var v interface{}
			err = red.XREADGROUP{
				Group: group,
				Consumer: "go",
				Count: 10,
				Block: time.Second*2,
				Streams: []red.QueryStream{
					{Key: key, ID: ">"},
				},
			}.Do(ctx, radixClient, v) ; if err != nil {
			return err
		}
			log.Print(v)
		}
	}
}