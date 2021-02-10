package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetExpireAt (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//   SET example_set_expire_at nimoc PXAT timestamp-milliseconds
	err = red.SET{
		Key: "example_set_expire_at",
		Value: "nimoc",
		// 只需配置 ExpireAt 为 time.Time，无需配置 EXAT PXAT ,goclub/redis 会将 ExpireAt 自动转换为 PXAT
		ExpireAt: time.Now().Add(time.Second),
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_expire_at"}.Do(ctx, radixClient))
	// "nimoc" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_expire_at"}.Do(ctx, radixClient))
	// "" false nil
}
