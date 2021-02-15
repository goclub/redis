package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHSET_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_hset"
	_=key
	{
		_,err := red.HSET{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_FORGET_ARGS) HSET Key is empty")
	}
	{
		_, err := red.HSET{Key: key, Field:"field", Value:"value"}.Do(ctx, Test{t, "HSET test_hset field value"})
		assert.NoError(t, err)
	}
	{
		_, err := red.HSET{Key: key, Multiple: []red.FieldValue{
			{"name", "nimoc"},
			{"age", "18"},
		}}.Do(ctx, Test{t, "HSET test_hset name nimoc age 18"})
		assert.NoError(t, err)
	}
	// 准备数据
	{
		_, err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
	}
	{
		added, err := red.HSET{Key:key, Field:"name",Value:"nimo1"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, added, uint(1))
	}
	{
		value, has,  err := red.HGET{
			Key: key,
			Field: "name",
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, value, "nimo1")
		assert.Equal(t, has, true)
	}
	{
		added, err := red.HSET{Key:key, Field:"name",Value:"nimo2"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, added, uint(0))
	}
	{
		value,has,  err := red.HGET{
			Key: key,
			Field: "name",
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, value, "nimo2")
		assert.Equal(t, has, true)
	}
	// 清空数据
	{
		_, err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
	}
	{
		added, err := red.HSET{
			Key:key,
			Multiple:[]red.FieldValue{
				{"name", "nimoc"},
				{"age", "18"},
			},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, added, uint(2))
	}
	{
		value,has,  err := red.HGET{
			Key: key,
			Field: "name",
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, value, "nimoc")
		assert.Equal(t, has, true)
	}
	{
		value, has, err := red.HGET{
			Key: key,
			Field: "age",
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, value, "18")
		assert.Equal(t, has, true)
	}
	{
		value, has, err := red.HGET{
			Key: key,
			Field: "invalid_field",
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, value, "")
		assert.Equal(t, has, false)
	}
	{
		values, err := red.HMGET{
			Key: key,
			Fields: []string{"name", "age", "invalid_field"},
		}.Do(ctx, radixClient)
		assert.NoError(t ,err)
		assert.Equal(t, values, []string{"nimoc", "18", ""})
	}
}