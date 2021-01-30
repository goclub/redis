package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)
func TestLRANGE_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_lrange"
	_=key
	{
		_,_ ,err := red.LRANGE{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_EMPTY_KEY)  LRANGE  key is empty")
	}
	{
		_,_ ,err := red.LRANGE{Key: key}.Do(ctx, Test{t, "LRANGE test_list_lrange 0 0"})
		assert.NoError(t, err)
	}
	{
		_,_ ,err := red.LRANGE{key, 1, -1}.Do(ctx, Test{t, "LRANGE test_list_lrange 1 -1"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		_, err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err := red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{})
		assert.Equal(t, isEmpty, true)
		// LPUSH key a
		_, err = red.LPUSH{Key: key, Value: "a"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{"a"})
		assert.Equal(t, isEmpty, false)
		// LPUSH key b
		_, err = red.LPUSH{Key: key, Value: "b"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{"b", "a"})
		assert.Equal(t, isEmpty, false)
		// LRANGE key -1 -1
		list, isEmpty, err = red.LRANGE{key, -1, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{"a"})
		assert.Equal(t, isEmpty, false)
		// LRANGE key -1 -1
		list, isEmpty, err = red.LRANGE{key, 0, 0}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{"b"})
		assert.Equal(t, isEmpty, false)
		// LPUSH key c b
		_, err = red.LPUSH{Key: key, Values: []string{"c", "d"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{"d", "c", "b", "a"})
		assert.Equal(t, isEmpty, false)
		// LRANGE key 10 10
		list, isEmpty, err = red.LRANGE{key, 10, 10}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, list, []string{})
		assert.Equal(t, isEmpty, true)
	}
}
func TestLPUSH_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_lpush"
	{
		_, err := red.LPUSH{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis(ERR_EMPTY_KEY)  LPUSH  key is empty")
	}
	{
		_, err := red.LPUSH{Key: key, Value: "a"}.Do(ctx, Test{t, "LPUSH test_list_lpush a"})
		assert.NoError(t, err)
	}
	{
		_, err := red.LPUSH{Key: key, Values: []string{"a", "b"}}.Do(ctx, Test{t, "LPUSH test_list_lpush a b"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		{
			_, err := red.DEL{Key: key}.Do(ctx, radixClient)
			assert.NoError(t, err)
		}
		// LPUSH key a
		{
			length, err := red.LPUSH{Key: key, Value: "a"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(1))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Do(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a"})
		}
		// LPUSH key b c
		{
			length, err := red.LPUSH{Key: key, Values: []string{"b", "c"}}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(3))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Do(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"c", "b", "a"})
		}

	}
}

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