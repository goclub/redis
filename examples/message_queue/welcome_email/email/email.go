package examplesMQEmail

import (
	"context"
	red "github.com/goclub/redis"
	exmapleMQ "github.com/goclub/redis/examples/message_queue/welcome_email/mesasge_queue"
	"github.com/mediocregopher/radix/v4"
	"log"
	"time"
)
func SendEmail(name string) {
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
			var list []radix.StreamEntries
			err = red.XREADGROUP{
				Group: group,
				Consumer: "go",
				Count: 100,
				Block: time.Second*1,
				Streams: []red.QueryStream{
					{Key: key, ID: ">"},
				},
			}.Do(ctx, radixClient, &list) ; if err != nil {
				return err
			}
			for _, item := range list {
				if item.Stream == key {
					for _, entry := range item.Entries {
						data := fieldsToMap(entry.Fields)
						name, hasName := data["name"]
						if hasName {
							SendEmail(name)
						} else {
							log.Print("can not found name", data, item.Stream, entry.ID)
						}
					}
				}
			}
		}
	}
}
func fieldsToMap(fields [][2]string) map[string]string {
	data := map[string]string{}
	for _, item := range fields {
		data[item[0]] = item[1]
	}
	return data
}