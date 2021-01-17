package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	radixClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
		panic(err)
	}
	client := red.NewRadix4Client(radixClient)

	_, err = red.SET{
		Key: "name",
		Value: "tim",
		Expires: time.Minute,
	}.Do(ctx, client) ; if err != nil {
		panic(err)
	}
	value, hasValue, err := red.GET{
		Key: "name",
	}.Do(ctx, client) ; if err != nil {
		panic(err)
	}
	log.Print(value, hasValue)
}