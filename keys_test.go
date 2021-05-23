package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	for _, client := range Connecters {
		redisCopy(t, client)
	}
}
func redisCopy(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "copy"
	destKey := "copy_dest"
	{ // DEL key
		_, err := DEL{Keys: []string{key, destKey}}.Do(ctx, client) ; assert.NoError(t, err)
	}
	{
		_, _, err := SET{NeverExpire: true, Key: key, Value: "v1"}.Do(ctx, client)  ; assert.NoError(t, err)
		reply, err := COPY{
			Source: key,
			Destination: destKey,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,reply, int64(1))
	}
	{
		// 已经存在的 dest key 不会被修改
		_, _, err := SET{NeverExpire: true, Key: key, Value: "v1"}.Do(ctx, client)  ; assert.NoError(t, err)
		reply, err := COPY{
			Source: key,
			Destination: destKey,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,reply, int64(0))
	}
}
func TestDump(t *testing.T) {
	for _, client := range Connecters {
		redisDump(t, client)
	}
}
func redisDump(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "dump"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	_, _, err = SET{NeverExpire: true, Key:key, Value:"Hello world"}.Do(ctx, client) ; assert.NoError(t, err)
	value, err := DUMP{Key:key}.Do(ctx, client)
	assert.Equal(t, value, "\x00\vHello world\t\x00\xcb!m\xae\x92ef\xe8")
}
func TestExists(t *testing.T) {
	for _, client := range Connecters {
		redisExists(t, client)
	}
}
func redisExists(t *testing.T, client Connecter) {
	ctx := context.TODO()
	keys := []string{"exist1", "exist2"}
	_, err := DEL{Keys: keys}.Do(ctx, client) ; assert.NoError(t, err)
	existsCount, err := EXISTS{
		Keys: keys,
	}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t,existsCount, uint64(0))
	_, _, err = SET{NeverExpire: true, Key: keys[0], Value: "a"}.Do(ctx, client) ; assert.NoError(t, err)
	{
		existsCount, err := EXISTS{
			Keys: keys,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,existsCount, uint64(1))
	}
	_, _, err = SET{NeverExpire: true, Key: keys[1], Value: "b"}.Do(ctx, client) ; assert.NoError(t, err)
	{
		existsCount, err := EXISTS{
			Keys: keys,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,existsCount, uint64(2))
	}
	_, err = DEL{Key: keys[0]}.Do(ctx, client) ; assert.NoError(t, err)
	{
		existsCount, err := EXISTS{
			Keys: keys,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,existsCount, uint64(1))
	}
}

func TestDel(t *testing.T) {
	for _, client := range Connecters {
		redisDel(t, client)
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
		assert.Equal(t, delCount, uint64(1))
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
		assert.Equal(t, delCount, uint64(2))
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


func TestPExpire(t *testing.T) {
	for _, client := range Connecters {
		redisPExpire(t, client)
	}
}
func redisPExpire(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "pexpire"
	_, err := DEL{Key:key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		reply, err := PEXPIRE{Key: key, Duration: time.Second}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,reply, int64(0))
		result, err := PTTL{
			Key: key,
		}.Do(ctx, client)  ; assert.NoError(t, err)
		assert.Equal(t,result.KeyDoesNotExist, true)
	}
	{
		_, _, err = SET{NeverExpire: true, Key:key, Value: "a"}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := PEXPIRE{Key: key, Duration: time.Second}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,reply, int64(1))
		result, err := PTTL{
			Key: key,
		}.Do(ctx, client)  ; assert.NoError(t, err)
		assert.Equal(t,result.TTL.Milliseconds() > 900, true)
		assert.Equal(t,result.TTL.Milliseconds() <= 1000, true)
	}
}

func TestKeys(t *testing.T) {
	for _, client := range Connecters {
		redisKeys(t, client)
	}
}
func redisKeys(t *testing.T, client Connecter) {
	ctx := context.TODO()
	keys := []string{"keys1", "keys2"}
	_, err := DEL{Keys:keys}.Do(ctx, client) ; assert.NoError(t, err)
	err = MSET{KeysValues: []KeyValue{
		{
			"keys1","v1",
		},
		{
			"keys2","v2",
		},
	}}.Do(ctx, client) ; assert.NoError(t, err)
	reply, err := KEYS{
		Pattern: "keys?",
	}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, len(reply), 2)
	findKey1 := false
	findKey2 := false
	for _, key := range reply {
		if key == "keys1" {
			findKey1 = true
		}
		if key == "keys2" {
			findKey2 = true
		}
	}
	assert.Equal(t, findKey1, true)
	assert.Equal(t, findKey2, true)
	{
		reply, err := KEYS{
			Pattern: "inexistence?",
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t,reply, []string(nil))
	}
}