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
		assert.EqualError(t, err, "goclub/redis: GET{} Key cannot be empty")
	}
	{// GET empty key
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "DEL", name))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.NoError(t, err)
	}
	{// GET valid key
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "SET", name, "abc"))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "abc")
		assert.Equal(t, hasValue, true)
		assert.NoError(t, err)
	}
	{// GET invalid valid key
		listKey := "test_get_list"
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "LPUSH", listKey, "mysql", "mongodb"))
		assert.NoError(t, err)
		value,hasValue, err := red.GET{Key:listKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.EqualError(t, err, "WRONGTYPE Operation against a key holding the wrong kind of value")
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
		assert.EqualError(t, err, "goclub/redis: DECR{} Key cannot be empty")
	}
	{// DECR empty key
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "DEL", name))
		assert.NoError(t, err)
		value, err := red.DECR{Key:name}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(-1))
		assert.NoError(t, err)
	}
	{// DECR valid key
		validKey := "test_valid_decr"
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "SET", validKey, "100"))
		assert.NoError(t, err)
		value, err := red.DECR{Key:validKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(99))
		assert.NoError(t, err)
	}
	{// DECR invalid key
		invalidKey := "test_invalid_decr"
		err := radixClient.Do(context.TODO(), radix.Cmd(nil, "SET", invalidKey, "234293482390480948029348230948"))
		assert.NoError(t, err)
		value, err := red.DECR{Key:invalidKey}.Do(context.TODO(), radixClient)
		assert.Equal(t, value, int64(0))
		assert.EqualError(t, err, "ERR value is not an integer or out of range")
	}
}