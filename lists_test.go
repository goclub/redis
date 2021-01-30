package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestBRPOPLPUSH_Do(t *testing.T) {
	ctx := context.Background()
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
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_EMPTY_KEY) BRPOPLPUSH Destination key is empty")
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
		}.Do(context.TODO(), Test{t, "BRPOPLPUSH src dest 0"})
		assert.NoError(t, err)
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
			Timeout: time.Second,
		}.Do(context.TODO(), Test{t, "BRPOPLPUSH src dest 1"})
		assert.NoError(t, err)
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
			Timeout: time.Minute,
		}.Do(context.TODO(), Test{t, "BRPOPLPUSH src dest 60"})
		assert.NoError(t, err)
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
			Timeout: time.Millisecond,
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis:(ERR_TIMEOUT) BRPOPLPUSH Timeout can not less at time.Second")
	}

	{
		sourceKey := "test_list_src"
		destinationKey := "test_list_desc"
		// 准备数据 [a] [b]
		{
			_, err := red.DEL{Keys:[]string{sourceKey, destinationKey}}.Do(ctx, radixClient)
			assert.NoError(t, err)
			_, err = red.Do(ctx, radixClient, nil, []string{"LPUSH", sourceKey, "a"})
			assert.NoError(t, err)
			_, err = red.Do(ctx, radixClient, nil, []string{"LPUSH",  destinationKey, "b"})
			assert.NoError(t, err)
		}
		// [a] -> [b] == [] [ab]
		{
			value, isNil, err := red.BRPOPLPUSH{
				Source: sourceKey,
				Destination: destinationKey,
			}.Do(ctx, radixClient)
			assert.NoError(t, err)
			assert.Equal(t, isNil, false)
			assert.Equal(t, value, "a")
			// 检查 src []
			{
				var list []string
				_, err := red.Do(ctx, radixClient, &list, []string{"LRANGE", sourceKey, "0", "-1"})
				assert.NoError(t, err)
				assert.Equal(t, list, []string{})
			}
			// 检查 desc [a b]
			{
				var list []string
				_, err := red.Do(ctx, radixClient, &list, []string{"LRANGE", destinationKey, "0", "-1"})
				assert.NoError(t, err)
				assert.Equal(t, list, []string{"a", "b"})
			}
		}
		// [a b] -> [b a]
		{
			value, isNil, err := red.BRPOPLPUSH{
				Source: destinationKey,
				Destination:destinationKey,
			}.Do(ctx, radixClient)
			assert.NoError(t, err)
			assert.Equal(t, isNil, false)
			assert.Equal(t, value, "b")
			// [b a]
			{
				var list []string
				_, err := red.Do(ctx, radixClient, &list, []string{"LRANGE", destinationKey, "0", "-1"})
				assert.NoError(t, err)
				assert.Equal(t, list, []string{"b", "a" })
			}
		}
		// timeout
		{
			// 准备数据
			{
				_, err := red.DEL{Key:"test_list_emptyKey"}.Do(ctx, radixClient)
				assert.NoError(t, err)
			}
			startTime := time.Now()
			value, isNil, err := red.BRPOPLPUSH{
				Source: "test_list_emptyKey",
				Destination: "test_list_emptyKey",
				Timeout: time.Second,
			}.Do(ctx, radixClient)
			assert.NoError(t, err)
			assert.Equal(t, isNil, true)
			assert.Equal(t, value, "")
			duration := time.Now().Sub(startTime)
			log.Print(duration.String())
			assert.Greater(t, duration, time.Second)
			assert.Less(t, duration, time.Millisecond * 1100)
		}
	}
}