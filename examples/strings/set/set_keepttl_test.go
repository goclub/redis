package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetKEEPTTL (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//
	err = red.SET{
		Key: "example_set_keep_ttl",
		Value: "x",
		Expire: time.Second * 1,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "x" true nil
	radixClient.DebugOnce()
	// SET example_set_keep_ttl xyz KEEPTTL
	err = red.SET{
		Key: "example_set_keep_ttl",
		Value: "y",
		// 注意 KEEPTTL 是 6.0 才支持的功能
		KeepTTL:true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "y" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "" false nil
}
