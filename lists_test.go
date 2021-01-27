package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBRPOPLPUSH_Do(t *testing.T) {
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "",
			Destination: "dest",
			Timeout: 1,
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_EMPTY_KEY) BRPOPLPUSH Source key is empty")
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "",
			Timeout: 1,
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_EMPTY_KEY) BRPOPLPUSH Destination key is empty")
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
			Timeout: 1,
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis:(ERR_DURATION) BRPOPLPUSH Timeout  can not set 1ns, maybe you forget time.Millisecond or time.time.Second")
	}
	// [a] -> [b] == [] [ab]
	{
		ctx := context.TODO()
		sourceKey := "test_list_src"
		destinationKey := "test_list_desc"
		_, err := red.DEL{Keys: []string{sourceKey, destinationKey}}.Do(ctx, radixClient) ; if err != nil {panic(err)}
		_, err = red.Do(context.TODO(), radixClient, nil, []string{"LPUSH", sourceKey, "a"}) ; if err != nil {
				panic(err)
			}
		_, err = red.Do(context.TODO(), radixClient, nil, []string{"LPUSH", destinationKey, "b"}) ; if err != nil {
			panic(err)
		}
		value, isNil, err := red.BRPOPLPUSH{
			Source: "test_list_src",
			Destination: "test_list_desc",
			Timeout: time.Second*1,
		}.Do(context.TODO(), radixClient)
		if err != nil {panic(err)}
		assert.Equal(t, value, "a")
		assert.Equal(t, isNil, false)
		arr := []string{}
		_, err = red.Do(ctx, radixClient, &arr, []string{"LRANGE", sourceKey, "0", "-1"}) ; if err != nil {panic(err)}
		assert.Equal(t, arr, []string{})
		_, err = red.Do(ctx, radixClient, &arr, []string{"LRANGE", destinationKey, "0", "-1"}) ; if err != nil {panic(err)}
		assert.Equal(t, arr, []string{"a", "b"})
	}
}