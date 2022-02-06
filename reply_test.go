package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReply(t *testing.T) {
	for _, client := range Connecters {
		redisReply(t, client)
	}
}

func redisReply(t *testing.T, client Connecter) {
	ctx := context.Background()
	// string
	{
		reply, isNil, err := client.Eval(ctx, Script{
			Script: `
			redis.call("SET", "reply_string", "1")
			return redis.call("GET", "reply_string")
		`,
		}) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		{
			v, err := reply.String() ; assert.NoError(t, err)
			assert.Equal(t, v, "1")
		}
		{
			v, err := reply.Int64() ; assert.NoError(t, err)
			assert.Equal(t, v, int64(1))
		}
		{
			v, err := reply.Uint64() ; assert.NoError(t, err)
			assert.Equal(t, v, uint64(1))
		}
		{
			v, err := reply.Int64Slice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(string) value(1) convert []OptionInt64")
			assert.Nil(t, v)
		}
		{
			v, err := reply.StringSlice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(string) value(1) convert []OptionString")
			assert.Nil(t, v)
		}
	}
	// string a
	{
		reply, isNil, err := client.Eval(ctx, Script{
			Script: `
			redis.call("SET", "reply_string_a", "a")
			return redis.call("GET", "reply_string_a")
		`,
		}) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		{
			v, err := reply.String() ; assert.NoError(t, err)
			assert.Equal(t, v, "a")
		}
		{
			v, err := reply.Int64()
			assert.EqualError(t, err, `strconv.ParseInt: parsing "a": invalid syntax`)
			assert.Equal(t, v, int64(0))
		}
		{
			v, err := reply.Int64Slice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(string) value(a) convert []OptionInt64")
			assert.Nil(t, v)
		}
		{
			v, err := reply.StringSlice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(string) value(a) convert []OptionString")
			assert.Nil(t, v)
		}
	}
	// int64
	{
		reply, isNil, err := client.Eval(ctx, Script{
			Script: `
			redis.call("SET", "reply_int64", "-2")
			return tonumber(redis.call("GET", "reply_int64"))
		`,
		}) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		{
			v, err := reply.String() ; assert.NoError(t, err)
			assert.Equal(t, v, "-2")
		}
		{
			v, err := reply.Int64() ; assert.NoError(t, err)
			assert.Equal(t, v, int64(-2))
		}
		{
			v, err := reply.Uint64()
			assert.Equal(t, v, uint64(0))
			assert.EqualError(t, err, "goclub/redis: -2 can not convert to uint64")
		}
		{
			v, err := reply.Int64Slice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(int64) value(-2) convert []OptionInt64")
			assert.Nil(t, v)
		}
		{
			v, err := reply.StringSlice()
			assert.EqualError(t, err, "goclub/redis: unexpected type(int64) value(-2) convert []OptionString")
			assert.Nil(t, v)
		}
	}
	// Int64Slice
	{
		reply, isNil, err := client.Eval(ctx, Script{
			Script: `
			redis.call("SADD", "reply_int64_slice", "1","2")
			return redis.call("SMEMBERS", "reply_int64_slice")
		`,
		}) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		{
			v, err := reply.String()
			assert.EqualError(t, err, "goclub/redis: unexpected type([]interface {}) value([1 2]) convert string")
			assert.Equal(t, v, "")
		}
		{
			v, err := reply.Int64()
			assert.EqualError(t, err, "goclub/redis: unexpected type([]interface {}) value([1 2]) convert int64")
			assert.Equal(t, v, int64(0))
		}
		{
			v, err := reply.Int64Slice()
			assert.NoError(t, err)
			assert.Equal(t,v, []OptionInt64{NewOptionInt64(1), NewOptionInt64(2)})
		}
		{
			v, err := reply.StringSlice()
			assert.NoError(t, err)
			assert.Equal(t,v, []OptionString{NewOptionString("1"), NewOptionString("2")})
		}
	}
}