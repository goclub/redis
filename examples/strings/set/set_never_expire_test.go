package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSetNeverExpire (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}

	//  SET example_set_never_expire hello
	err = red.SET{
		Key: "example_set_never_expire",
		Value: "hello",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	log.Print(red.GET{Key: "example_set_never_expire"}.Do(ctx, radixClient))
	// "hello" true nil

	// 如果你没有传入 NeverExpire: true ，则会返回一个错误提醒你可能忘记配置过期时间
	err = red.SET{
		Key:"example_set_never_expire",
		Value: "hello",
	}.Do(ctx, radixClient)
	assert.Error(t, err, red.ErrSetForgetTimeToLive)
}
