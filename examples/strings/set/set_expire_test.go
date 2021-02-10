package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetExpire (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//  SET example_set abc PX 1000
	err = red.SET{
		Key: "example_set_expire",
		Value: "abc",
		// 只需配置 Expire 为 time.Duration，无需配置 EX PX ,goclub/redis 会将 Expire 自动转换为 PX
		Expire: time.Second * 1,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_expire"}.Do(ctx, radixClient))
	// "abc" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_expire"}.Do(ctx, radixClient))
	// "" false nil
}
