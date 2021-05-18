package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func redisAppend(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "append"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		length, err := APPEND{
			Key: key,
			Value: "a",
		}.Do(ctx, client) ; assert.NoError(t ,err)
		assert.Equal(t, length, uint64(1))
		value, hasValue, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, hasValue, true)
		assert.Equal(t, value, "a")
	}
	{
		length, err := APPEND{
			Key: key,
			Value: "b",
		}.Do(ctx, client) ; assert.NoError(t ,err)
		assert.Equal(t, length, uint64(2))
		value, hasValue, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, hasValue, true)
		assert.Equal(t, value, "ab")
	}
	{
		length, err := APPEND{
			Key: key,
			Value: "cd",
		}.Do(ctx, client) ; assert.NoError(t ,err)
		assert.Equal(t, length, uint64(4))
		value, hasValue, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, hasValue, true)
		assert.Equal(t, value, "abcd")
	}
}


func TestAppend(t *testing.T) {
	for _, client := range Connecters {
		redisAppend(t, client)
	}
}

func redisDel(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "del"
	key2 := "del2"
	{// SET key v
		_, _, err := client.DoStringReply(ctx, []string{"SET", key, "v"}) ; assert.NoError(t, err)
	}
	{// DEL key
		delCount, err := DEL{Key: key}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(1))
		reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key})
		assert.NoError(t, err)
		assert.Equal(t, reply, "")
		assert.Equal(t, isNil, true)}
	{
		// MSET key "nimo" key2 "nico"
		_, _, err := client.DoStringReply(ctx, []string{"MSET", key, "nimo", key2, "nico"}) ; assert.NoError(t, err)
	}
	{// DEL key key2
		delCount, err := DEL{Keys: []string{key, key2}}.Do(ctx, client)
		assert.Equal(t, delCount, uint(2))
		assert.NoError(t, err)
		// GET key
		{
			reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key})
			assert.Equal(t, reply, "")
			assert.Equal(t, isNil, true)
			assert.NoError(t, err)
		}
		// GET key2
		{
			reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key2})
			assert.Equal(t, reply, "")
			assert.Equal(t, isNil, true)
			assert.NoError(t, err)
		}
	}


}
func TestDel(t *testing.T) {
	for _, client := range Connecters {
		redisDel(t, client)
	}
}

func redisPTTL(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "pttl"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{ // PTTL key
		result, err := PTTL{
			Key: key,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, result, ResultTTL{
			TTL: 0,
			KeyDoesNotExist: true,
			NeverExpire: false,
		})
	}
	{ // SET key value
		_, _, err := client.DoStringReply(ctx, []string{"SET", key, "goclub"}) ;assert.NoError(t, err)
	}
	{ // PTTL key
		result, err := PTTL{
			Key: key,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, result, ResultTTL{
			TTL: 0,
			KeyDoesNotExist: false,
			NeverExpire: true,
		})
	}
	{ // SET key value EX 2
		_, _, err := client.DoStringReply(ctx, []string{"SET", key, "goclub", "EX", "2"}) ;assert.NoError(t, err)
	}
	{ // PTTL key
		result, err := PTTL{
			Key: key,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, result.NeverExpire, false)
		assert.Equal(t, result.KeyDoesNotExist, false)
		assert.Equal(t, result.TTL > time.Second, true)
		assert.Equal(t, result.TTL < time.Second*2, true)
	}
}

func TestPTTL(t *testing.T) {
	for _, client := range Connecters {
		redisPTTL(t, client)
	}
}
func redisSet(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "set"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{// GET key
		value, hasValue, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.NoError(t, err)
	}
	{// SET key value
		isNil, err := SET{
			Key:         key,
			Value:       "goclub",
			NeverExpire: true,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
	}
	{// GET key
		value, hasValue, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "goclub")
		assert.Equal(t, hasValue, true)
		assert.NoError(t, err)
	}
	// @NEXT 写上 SET 的测试

}
func TestSet(t *testing.T) {
	for _, client := range Connecters {
		redisSet(t, client)
	}
}

func redisGet(t *testing.T,  client Connecter) {
	ctx := context.TODO()
	key := "get"
	// DEL key
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{// GET key
		value, hasValue, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "")
		assert.Equal(t, hasValue, false)
		assert.NoError(t, err)
	}
	{// SET key value
		newValue := "nimo" + strconv.Itoa(time.Now().Nanosecond())
		_, _, err = client.DoStringReply(ctx, []string{"SET", key, newValue}) ; assert.NoError(t, err)
		_="GET key"
		value, hasValue, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, newValue)
		assert.Equal(t, hasValue, true)
		assert.NoError(t, err)
	}

}

func TestGet(t *testing.T) {
	for _, client := range Connecters {
		redisGet(t, client)
	}
}
