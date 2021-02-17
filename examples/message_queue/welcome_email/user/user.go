package examplesMQUser

import (
	"context"
	red "github.com/goclub/redis"
	exmapleMQ "github.com/goclub/redis/examples/message_queue/welcome_email/mesasge_queue"
	"log"
)




func PublishUserSignInMessage(ctx context.Context, radixClient red.DriverRadixClient4, name string) (err error) {
	mqdata := exmapleMQ.UserSignInMessage{
		Name: name,
	}
	data, err  := red.StructFieldValues(mqdata) ; if err != nil {
		return
	}
	streamID, err := red.XADD{
		Key: mqdata.StreamKey(),
		FieldValues:data,
	}.Do(ctx, radixClient) ; if err != nil {
		return
	}
	log.Print("streamID:", streamID.String())
	return
}