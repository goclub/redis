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
		_ ,_,err := red.LRANGE{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) LRANGE Key can not be empty")
	}
	{
		_, _ ,err := red.LRANGE{Key: key}.Do(ctx, Test{t, "LRANGE test_list_lrange 0 0"})
		assert.NoError(t, err)
	}
	{
		_, _ ,err := red.LRANGE{key, 1, -1}.Do(ctx, Test{t, "LRANGE test_list_lrange 1 -1"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		_, err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list,isEmpty, err := red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, true)
		assert.Equal(t, list, []string(nil))
		// LPUSH key a
		_, err = red.LPUSH{Key: key, Value: "a"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list,isEmpty,  err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, false)
		assert.Equal(t, list, []string{"a"})
		
		// LPUSH key b
		_, err = red.LPUSH{Key: key, Value: "b"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, false)
		assert.Equal(t, list, []string{"b", "a"})
		
		// LRANGE key -1 -1
		list, isEmpty, err = red.LRANGE{key, -1, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, false)
		assert.Equal(t, list, []string{"a"})
		
		// LRANGE key -1 -1
		list, isEmpty, err = red.LRANGE{key, 0, 0}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, false)
		assert.Equal(t, list, []string{"b"})
		
		// LPUSH key c b
		_, err = red.LPUSH{Key: key, Values: []string{"c", "d"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LRANGE key 0 -1
		list, isEmpty, err = red.LRANGE{key, 0, -1}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, false)
		assert.Equal(t, list, []string{"d", "c", "b", "a"})
		
		// LRANGE key 10 10
		list, isEmpty, err = red.LRANGE{key, 10, 10}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, true)
		assert.Equal(t, list, []string(nil))
		
		// LRANGE key 10 4
		list,isEmpty, err = red.LRANGE{key, 10, 4}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isEmpty, true)
		assert.Equal(t, list, []string(nil))
		
	}
}
func TestLPUSH_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_lpush"
	{
		_, err := red.LPUSH{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis: (ERR_FORGET_ARGS) LPUSH Key can not be empty")
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
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
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
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"c", "b", "a"})
		}

	}
}
func TestLPUSHX_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_lpushx"
	{
		_, err := red.LPUSHX{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis: (ERR_FORGET_ARGS) LPUSHX Key can not be empty")
	}
	{
		_, err := red.LPUSHX{Key: key, Value: "a"}.Do(ctx, Test{t, "LPUSHX test_list_lpushx a"})
		assert.NoError(t, err)
	}
	{
		_, err := red.LPUSHX{Key: key, Values: []string{"a", "b"}}.Do(ctx, Test{t, "LPUSHX test_list_lpushx a b"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		{
			_, err := red.DEL{Key: key}.Do(ctx, radixClient)
			assert.NoError(t, err)
		}
		// LPUSHX key a
		{
			length, err := red.LPUSHX{Key: key, Value: "a"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(0))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{})
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
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a"})
		}
		// LPUSHX key b
		{
			length, err := red.LPUSHX{Key: key, Value: "b"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(2))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"b","a"})
		}

	}
}
func TestRPUSH_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_rpush"
	{
		_, err := red.RPUSH{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis: (ERR_FORGET_ARGS) RPUSH Key can not be empty")
	}
	{
		_, err := red.RPUSH{Key: key, Value: "a"}.Do(ctx, Test{t, "RPUSH test_list_rpush a"})
		assert.NoError(t, err)
	}
	{
		_, err := red.RPUSH{Key: key, Values: []string{"a", "b"}}.Do(ctx, Test{t, "RPUSH test_list_rpush a b"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		{
			_, err := red.DEL{Key: key}.Do(ctx, radixClient)
			assert.NoError(t, err)
		}
		// RPUSH key a
		{
			length, err := red.RPUSH{Key: key, Value: "a"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(1))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a"})
		}
		// RPUSH key b c
		{
			length, err := red.RPUSH{Key: key, Values: []string{"b", "c"}}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(3))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a", "b", "c"})
		}

	}
}
func TestRPUSHX_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_rpushx"
	{
		_, err := red.RPUSHX{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis: (ERR_FORGET_ARGS) RPUSHX Key can not be empty")
	}
	{
		_, err := red.RPUSHX{Key: key, Value: "a"}.Do(ctx, Test{t, "RPUSHX test_list_rpushx a"})
		assert.NoError(t, err)
	}
	{
		_, err := red.RPUSHX{Key: key, Values: []string{"a", "b"}}.Do(ctx, Test{t, "RPUSHX test_list_rpushx a b"})
		assert.NoError(t, err)
	}
	{
		// 准备数据
		{
			_, err := red.DEL{Key: key}.Do(ctx, radixClient)
			assert.NoError(t, err)
		}
		// RPUSHX key a
		{
			length, err := red.RPUSHX{Key: key, Value: "a"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(0))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{})
		}
		// RPUSH key a
		{
			length, err := red.RPUSH{Key: key, Value: "a"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(1))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a"})
		}
		// RPUSHX key b
		{
			length, err := red.RPUSHX{Key: key, Value: "b"}.Do(ctx, radixClient)
			assert.Equal(t, length, uint(2))
			assert.NoError(t, err)
		}
		// check
		{
			var list []string
			_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", key, "0", "-1"})
			assert.NoError(t, err)
			assert.Equal(t, list, []string{"a", "b"})
		}

	}
}
func TestLPOP_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_lpop"
	{
		_,isNil, err := red.LPOP{}.Do(ctx, Test{})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) LPOP Key can not be empty")
		assert.Equal(t, isNil, false)
	}
	{
		_,isNil, err := red.LPOP{Key: key}.Do(ctx, Test{t, "LPOP test_list_lpop"})
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
	}
	{
		// 准备数据
		_,err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		_, err = red.RPUSH{Key: key, Values: []string{"a", "b", "c"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LPOP key
		value, isNil, err := red.LPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "a")
		// LPOP key
		value, isNil, err = red.LPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "b")
		// LPOP key
		value, isNil, err = red.LPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "c")
		// LPOP key
		value, isNil, err = red.LPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, value, "")
	}
	{
		// 准备数据
		_,err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		_, err = red.RPUSH{Key: key, Values: []string{"a", "b", "c"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// LPOP key 0
		_, _, err = red.LPOPCount{Key: key}.Do(ctx, radixClient)
		assert.EqualError(t, err, "goclub/redis(ERR_COUNT_CAN_NOT_BE_ZERO) LPOPCount{}.Count can not be zero")
		// LPOP key 2
		list, isNil, err := red.LPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, list, []string{"a","b"})
		// LPOP key 2
		list, isNil, err = red.LPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, list, []string{"c"})
		// LPOP key 2
		list, isNil, err = red.LPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, list, []string(nil))
	}
}
func TestRPOP_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_list_rpop"
	{
		_,isNil, err := red.RPOP{}.Do(ctx, Test{})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) RPOP Key can not be empty")
		assert.Equal(t, isNil, false)
	}
	{
		_,isNil, err := red.RPOP{Key: key}.Do(ctx, Test{t, "RPOP test_list_rpop"})
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
	}
	{
		// 准备数据
		_,err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		_, err = red.RPUSH{Key: key, Values: []string{"a", "b", "c"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// RPOP key
		value, isNil, err := red.RPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "c")
		// RPOP key
		value, isNil, err = red.RPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "b")
		// RPOP key
		value, isNil, err = red.RPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "a")
		// RPOP key
		value, isNil, err = red.RPOP{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, value, "")
	}
	{
		// 准备数据
		_,err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		_, err = red.RPUSH{Key: key, Values: []string{"a", "b", "c"}}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// RPOP key 0
		_, _, err = red.RPOPCount{Key: key}.Do(ctx, radixClient)
		assert.EqualError(t, err, "goclub/redis(ERR_COUNT_CAN_NOT_BE_ZERO) RPOPCount{}.Count can not be zero")
		// RPOP key 2
		list, isNil, err := red.RPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, list, []string{"c","b"})
		// RPOP key 2
		list, isNil, err = red.RPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, list, []string{"a"})
		// RPOP key 2
		list, isNil, err = red.RPOPCount{Key:key, Count: 2}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, list, []string(nil))
	}
}

func TestBRPOPLPUSH_Do(t *testing.T) {
	ctx := context.Background()
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "",
			Destination: "dest",
			Timeout: red.Second{1},
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) BRPOPLPUSH Source can not be empty")
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "",
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) BRPOPLPUSH Destination can not be empty")
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
			Timeout: red.Second{1},
		}.Do(context.TODO(), Test{t, "BRPOPLPUSH src dest 1"})
		assert.NoError(t, err)
	}
	{
		_, _, err := red.BRPOPLPUSH{
			Source: "src",
			Destination: "dest",
			Timeout: red.Second{60},
		}.Do(context.TODO(), Test{t, "BRPOPLPUSH src dest 60"})
		assert.NoError(t, err)
	}

	{
		sourceKey := "test_list_src"
		destinationKey := "test_list_desc"
		// 准备数据 [a] [b]
		{
			_, err := red.DEL{Keys:[]string{sourceKey, destinationKey}}.Do(ctx, radixClient)
			assert.NoError(t, err)
			_, err = red.Command(ctx, radixClient, nil, []string{"LPUSH", sourceKey, "a"})
			assert.NoError(t, err)
			_, err = red.Command(ctx, radixClient, nil, []string{"LPUSH",  destinationKey, "b"})
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
				_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", sourceKey, "0", "-1"})
				assert.NoError(t, err)
				assert.Equal(t, list, []string{})
			}
			// 检查 desc [a b]
			{
				var list []string
				_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", destinationKey, "0", "-1"})
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
				_, err := red.Command(ctx, radixClient, &list, []string{"LRANGE", destinationKey, "0", "-1"})
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
				Timeout: red.Second{1},
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