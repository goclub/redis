package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestAppend(t *testing.T) {
	for _, client := range Connecters {
		redisAppend(t, client)
	}
}

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

func TestBitCount(t *testing.T) {
	for _, client := range Connecters {
		redisBitCount(t, client)
	}
}
func redisBitCount(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "bitcount"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	_, err = SET{Key: key, Value:"foobar", NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
	{
		length, err := BITCOUNT{
			Key: key,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(26))
	}
	{
		length, err := BITCOUNT{
			Key: key,
			Start: Uint64(0),
			End: Uint64(0),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(4))
	}
	{
		length, err := BITCOUNT{
			Key: key,
			Start: Uint64(1),
			End: Uint64(1),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(6))
	}
	{
		length, err := BITCOUNT{
			Key: key,
			Start: Uint64(1),
			End: Uint64(2),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(12))
	}
}

func TestBitField(t *testing.T) {
	for _, client := range Connecters {
		redisBitField(t, client)
	}
}
func redisBitField(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "bitfield"
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "i5", "100", "1",
				"GET", "u4", "0",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{1, 0})
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := BITFIELD{
			Key: key,
			Args: []string{
				"SET", "i8", "#0", "100",
				"SET", "i8", "#1", "200",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{0, 0})
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"GET", "i8", "#0",
				"GET", "i8", "#1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{100, -56})
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{1, 1})
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{2, 2})
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{3, 3})
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{0, 3})
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := BITFIELD{
			Key: key,
			Args: []string{
				"OVERFLOW", "FAIL",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, []int64{1})
	}

}

func TestBitop(t *testing.T) {
	for _, client := range Connecters {
		redisBitop(t, client)
	}
}
func redisBitop(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "BITOP"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	_, err = DEL{Key: "bittop_dest"}.Do(ctx, client) ; assert.NoError(t, err)
	_, err = SET{Key: key +"1", Value: "foobar",NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
	_, err = SET{Key: key +"2", Value: "abcdef",NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
	size, err := BITOP{
		AND: true,
		DestKey: "bittop_dest",
		Keys: []string{key+"1",key+"2"},
	}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, size, uint64(6))
	value, hasValue, err := GET{Key: "bittop_dest"}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, hasValue, true)
	assert.Equal(t, value, "`bc`ab")
}

func TestBitpos(t *testing.T) {
	for _, client := range Connecters {
		redisBitpos(t, client)
	}
}
func redisBitpos(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "bitpos"
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		_, err = SET{Key: key, Value: "\xff\xf0\x00",NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
		position, err := BITPOS{
			Key: key,
			Bit:0,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, position, int64(12))
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		_, err = SET{Key: key, Value: "\x00\xff\xf0",NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
		position, err := BITPOS{
			Key: key,
			Bit:1,
			Start: Uint64(0),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, position, int64(8))
		position, err = BITPOS{
			Key: key,
			Bit:1,
			Start: Uint64(2),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, position, int64(16))
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		_, err = SET{Key: key, Value: "\x00\x00\x00",NeverExpire: true}.Do(ctx, client) ; assert.NoError(t, err)
		position, err := BITPOS{
			Key: key,
			Bit:1,
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, position, int64(-1))
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

func TestPTTL(t *testing.T) {
	for _, client := range Connecters {
		redisPTTL(t, client)
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

func TestSet(t *testing.T) {
	for _, client := range Connecters {
		redisSet(t, client)
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


func TestGet(t *testing.T) {
	for _, client := range Connecters {
		redisGet(t, client)
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

