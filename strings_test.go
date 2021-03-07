package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGET_Do(t *testing.T) {
	name := "test_get"
	{
		_,_,_ = red.GET{
			Key: name,
		}.Do(context.TODO(), Test{t, "GET test_get"})
	}
	{
		_,_,err := red.GET{
			Key: "",
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis:  GET{} Key cannot be empty")
	}
	{// GET empty key
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "DEL", name))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.NoError(t, err)
	}
	{// GET valid key
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "SET", name, "abc"))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "abc")
		assert.Equal(t, hasValue, true)
		assert.NoError(t, err)
	}
	{// GET invalid valid key
		listKey := "test_get_list"
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "LPUSH", listKey, "mysql", "mongodb"))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:listKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.EqualError(t, err, "WRONGTYPE Operation against a key holding the wrong kind of value")
	}
}

func TestDEL_Do(t *testing.T) {
	name := "test_del"
	name2 := "test_del_2"
	{
		_, err := red.DEL{
			Keys: nil,
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis:  DEL{} Keys cannot be empty")
	}
	{
		_, err := red.DEL{
			Keys: []string{},
		}.Do(context.TODO(), Test{t, ""})
		assert.EqualError(t, err, "goclub/redis:  DEL{} Keys cannot be empty")
	}
	{
		_, _ = red.DEL{
			Keys: []string{name},
		}.Do(context.TODO(), Test{t, "DEL test_del"})
	}
	{
		_, _ = red.DEL{
			Keys: []string{name, name2,},
		}.Do(context.TODO(), Test{t, "DEL test_del test_del_2"})
	}
	{
		err := red.SET{
			Key: name,
			Value: "a",
			NeverExpire: true,
		}.Do(context.TODO(), radixClient)
		assert.NoError(t, err)
		delCount, err := red.DEL{
			Keys:[]string{name},
		}.Do(context.TODO(), radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(1))
	}
	{
		_, err := red.Command(context.TODO(), radixClient, nil, []string{"DEL", name2})
		assert.NoError(t, err)
		delCount, err := red.DEL{
			Keys:[]string{name2},
		}.Do(context.TODO(), radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(0))
	}
	{
		_, err := red.Command(context.TODO(), radixClient, nil, []string{"DEL", name, name2})
		assert.NoError(t, err)
		delCount, err := red.DEL{
			Keys:[]string{name, name2},
		}.Do(context.TODO(), radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(0))
	}
}
func TestDECR_Do(t *testing.T) {
	name := "test_decr"
	{
		_,_ = red.DECR{
			Key: name,
		}.Do(context.TODO(), Test{t, "DECR test_decr"})
	}
	{
		_,err := red.DECR{
			Key: "",
		}.Do(context.TODO(), Test{t, "DECR test_decr"})
		assert.EqualError(t, err, "goclub/redis:  DECR{} Key cannot be empty")
	}
	{// DECR empty key
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "DEL", name))
		assert.NoError(t, err)
		value, err := red.DECR{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(-1))
		assert.NoError(t, err)
	}
	{// DECR valid key
		validKey := "test_valid_decr"
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "SET", validKey, "100"))
		assert.NoError(t, err)
		value, err := red.DECR{Key:validKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(99))
		assert.NoError(t, err)
	}
	{// DECR invalid key
		invalidKey := "test_invalid_decr"
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "SET", invalidKey, "234293482390480948029348230948"))
		assert.NoError(t, err)
		value, err := red.DECR{Key:invalidKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(0))
		assert.EqualError(t, err, "ERR value is not an integer or out of range")
	}
}

func TestINCR_Do(t *testing.T) {
	name := "test_INCR"
	{
		_,_ = red.INCR{
			Key: name,
		}.Do(context.TODO(), Test{t, "INCR test_INCR"})
	}
	{
		_,err := red.INCR{
			Key: "",
		}.Do(context.TODO(), Test{t, "INCR test_INCR"})
		assert.EqualError(t, err, "goclub/redis:  INCR{} Key cannot be empty")
	}
	{// INCR empty key
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "DEL", name))
		assert.NoError(t, err)
		value, err := red.INCR{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(1))
		assert.NoError(t, err)
	}
	{// INCR valid key
		validKey := "test_valid_INCR"
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "SET", validKey, "100"))
		assert.NoError(t, err)
		value, err := red.INCR{Key:validKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(101))
		assert.NoError(t, err)
	}
	{// INCR invalid key
		invalidKey := "test_invalid_INCR"
		err := radixClient.Core.Do(context.TODO(), radix.Cmd(nil, "SET", invalidKey, "234293482390480948029348230948"))
		assert.NoError(t, err)
		value, err := red.INCR{Key:invalidKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(0))
		assert.EqualError(t, err, "ERR value is not an integer or out of range")
	}
}

func TestAPPEND_Do(t *testing.T) {
	name := "test_append"
	ctx := context.TODO()
	{
		_, err := red.APPEND{Key: "", Value: ""}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) APPEND Key can not be empty")
	}
	{
		_, err := red.APPEND{Key: name, Value: ""}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) APPEND Value can not be empty")
	}
	{
		_, err := red.APPEND{Key: name, Value: "1"}.Do(ctx, Test{t, "APPEND test_append 1"})
		assert.NoError(t, err)
	}
	{
		_, err := red.DEL{Key: name}.Do(ctx, radixClient)
		assert.NoError(t, err)
	}
	{
		length, err := red.APPEND{Key: name, Value: "1"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), length)
		value, hasValue, err := red.GET{Key:name}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, hasValue, true)
		assert.Equal(t, value, "1")
	}
	{
		length, err := red.APPEND{Key: name, Value: "2"}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, uint(2), length)
		value, hasValue, err := red.GET{Key:name}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, hasValue, true)
		assert.Equal(t, value, "12")
	}
}


func TestBIT(t *testing.T) {
	name := "test_bit"
	_=name
	ctx := context.TODO()
	{
		_, err := red.GETBIT{Key: ""}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) GETBIT Key can not be empty")
	}
	{
		_, err := red.GETBIT{Key: name}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) GETBIT offset can not be empty")
	}
	{
		_, err := red.GETBIT{Key: name, Offset: red.Uint32(0)}.Do(ctx, Test{t, "GETBIT test_bit 0"})
		assert.NoError(t, err)
	}

	{
		_, err := red.SETBIT{Key: ""}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) SETBIT Key can not be empty")
	}
	{
		_, err := red.SETBIT{Key: name}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) SETBIT offset can not be empty")
	}
	{
		_, err := red.SETBIT{Key: name, Offset: red.Uint32(0)}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: (ERR_FORGET_ARGS) SETBIT value can not be empty")
	}
	{
		_, err := red.SETBIT{Key: name, Offset: red.Uint32(0), Value: red.Uint(2)}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: SETBIT value must be 0 or 1, can not be 2")
	}
	{
		_, err := red.SETBIT{Key: name, Offset: red.Uint32(0), Value: red.Uint(1)}.Do(ctx, Test{t, "SETBIT test_bit 0 1"})
		assert.NoError(t, err)
	}
	
}