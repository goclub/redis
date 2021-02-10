package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetXX (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}

	//  SET example_set_xx hello XX
	ok, err := red.SETXX{
		Key: "example_set_xx",
		Value: "hello1",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第一次失败，因为 key 不存在
	assert.Equal(t, ok, false)

	// 第二次 SET key value 设置值
	//  SET example_set_xx hello
	err = red.SET{
		Key: "example_set_xx",
		Value: "hello2",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	//  SET example_set_xx hello XX
	ok, err = red.SETXX{
		Key: "example_set_xx",
		Value: "hello3",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第三次成功，因为 key 存在
	assert.Equal(t, ok, true)


	// GET example_set_xx
	value, hasValue, err := red.GET{Key:"example_set_xx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	assert.Equal(t, value, "hello3")
	assert.Equal(t, hasValue, true)
}
