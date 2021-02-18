package main

import (
	"context"
	"errors"
	red "github.com/goclub/redis"
	exampleMQData "github.com/goclub/redis/examples/message_queue/data"
	xtest "github.com/goclub/test"
	"github.com/mediocregopher/radix/v4"
	radix4 "github.com/redis-driver/mediocregopher-radix-v4"
	"log"
	"strings"
	"time"
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
	streamKey := exampleMQData.UserSignIn{}.StreamKey()
	group := "welcome_email"
	err := red.XGROUPCreate{
		Key: streamKey,
		Group: group,
		ID: "0",
		MKSTREAM: true,
	}.Do(ctx, client) ; if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			// 忽略
		} else {
			panic(err)
		}
	}
	for {
		var list []radix.StreamEntries
		err = red.XREADGROUP{
			Group: group,
			Consumer: "go-1",
			Count: 10,
			Block: time.Second,
			Streams: []red.QueryStream{
				{
					Key:streamKey,
					ID: ">",
				},
			},
		}.Do(ctx, client, &list) ; if err != nil {
			log.Print(err)
		}
		var entries  []radix.StreamEntry
		for _, item := range list {
			if item.Stream == streamKey {
				entries = item.Entries
			}
		}
		for _, entry := range entries {
			data := exampleMQData.UserSignIn{}
			err := red.StructScanByField(&data, radix4.StreamEntryFields(entry.Fields)) ; if err != nil {
				log.Print(err)
				continue
			}
			name := apiUserIDByUserName(data.UserID)
			time.Sleep(time.Second)
			log.Print("Send email: " , data.UserID, " Hi (" , name, "), welcome to goclub/redis.")
			ackCount, err := red.XACK{
				Key: streamKey,
				Group: group,
				StreamID: entry.ID.String(),
			}.Do(ctx, client) ;  if err != nil {
				log.Print(err)
				continue
			}
			if ackCount != 1 {
				log.Print(errors.New("stream:"  + streamKey + " group:" + group +" streamID:" + entry.ID.String() + ": ack count must be 1"))
			}
		}
	}
}
func apiUserIDByUserName(userID string) string {
	return xtest.Name()
}
