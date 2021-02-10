package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetNX (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	_, err = red.DEL{Key: "example_set_nx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	//  SET example_set_nx hello1 NX
	ok, err := red.SETNX{
		Key: "example_set_nx",
		Value: "hello1",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第一次成功
	assert.Equal(t, ok, true)

	//  SET example_set_nx hello2 NX
	ok, err = red.SETNX{
		Key: "example_set_nx",
		Value: "hello2",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第二次因为 key 已存在所以失败
	assert.Equal(t, ok, false)

	// GET example_set_nx
	value, hasValue, err := red.GET{Key:"example_set_nx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	assert.Equal(t, value, "hello1")
	assert.Equal(t, hasValue, true)
}
