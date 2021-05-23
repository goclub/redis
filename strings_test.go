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
		value, isNil, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "a")
	}
	{
		length, err := APPEND{
			Key: key,
			Value: "b",
		}.Do(ctx, client) ; assert.NoError(t ,err)
		assert.Equal(t, length, uint64(2))
		value, isNil, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "ab")
	}
	{
		length, err := APPEND{
			Key: key,
			Value: "cd",
		}.Do(ctx, client) ; assert.NoError(t ,err)
		assert.Equal(t, length, uint64(4))
		value, isNil, err := GET{Key:key}.Do(ctx, client) ; assert.NoError(t,err)
		assert.Equal(t, isNil, false)
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
			StartByte: NewOptionInt64(0),
			EndByte: NewOptionInt64(0),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(4))
	}
	{
		length, err := BITCOUNT{
			Key: key,
			StartByte: NewOptionInt64(1),
			EndByte: NewOptionInt64(1),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, length, uint64(6))
	}
	{
		length, err := BITCOUNT{
			Key: key,
			StartByte: NewOptionInt64(1),
			EndByte: NewOptionInt64(2),
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
	OptionInt64Slice := func (args ...int64) (list []OptionInt64) {
		for _, item := range args {
			list = append(list, NewOptionInt64(item))
		}
		return
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		reply, err := BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "i5", "100", "1",
				"GET", "u4", "0",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, OptionInt64Slice(1, 0))
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
		assert.Equal(t, reply, OptionInt64Slice(0,0))
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"GET", "i8", "#0",
				"GET", "i8", "#1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, OptionInt64Slice(100, -56))
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
		assert.Equal(t, reply, OptionInt64Slice(1, 1))
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, OptionInt64Slice(2, 2))
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, OptionInt64Slice(3, 3))
		reply, err = BITFIELD{
			Key: key,
			Args: []string{
				"INCRBY", "u2", "100", "1",
				"OVERFLOW", "SAT",
				"INCRBY", "u2", "102", "1",
			},
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, reply, OptionInt64Slice(0, 3))

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
		assert.Equal(t, reply, OptionInt64Slice(1))
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
	value, isNil, err := GET{Key: "bittop_dest"}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, isNil, false)
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
			Start: NewOptionUint64(0),
		}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, position, int64(8))
		position, err = BITPOS{
			Key: key,
			Bit:1,
			Start: NewOptionUint64(2),
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

func TestDecr(t *testing.T) {
	for _, client := range Connecters {
		redisDecr(t, client)
	}
}
func redisDecr(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "decr"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		newValue, err := DECR{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-1))
	}
	{
		newValue, err := DECR{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-2))
	}
}
func TestDecrBy(t *testing.T) {
	for _, client := range Connecters {
		redisDecrBy(t, client)
	}
}
func redisDecrBy(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "decrby"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		newValue, err := DECRBY{Key: key, Decrement: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-1))
	}
	{
		newValue, err := DECRBY{Key: key, Decrement: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-2))
	}
	{
		newValue, err := DECRBY{Key: key, Decrement: -1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-1))
	}
	{
		newValue, err := DECRBY{Key: key, Decrement: -1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(0))
	}
	{
		newValue, err := DECRBY{Key: key, Decrement: -2}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(2))
	}
	{
		newValue, err := DECRBY{Key: key, Decrement: 4}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(-2))
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
		assert.Equal(t, result.TTL <= time.Second*2, true)
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
		value, isNil, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "")
		assert.Equal(t, isNil, true)
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
		value, isNil, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "goclub")
		assert.Equal(t, isNil, false)
		assert.NoError(t, err)
	}
	// @NEXT 写上 SET 的测试

}

func TestSetBit(t *testing.T) {
	for _, client := range Connecters {
		redisSetBit(t, client)
	}
}
func redisSetBit(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "setbit"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		value, err := SETBIT{Key: key, Offset: 0, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))

		value, err = SETBIT{Key: key, Offset: 0, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))

		value, err = GETBIT{Key: key, Offset: 10}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))
	}
	{
		value, err := SETBIT{Key: key, Offset: 20, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))

		value, err = SETBIT{Key: key, Offset: 20, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))

		value, err = SETBIT{Key: key, Offset: 20, Value: 0}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))

		value, err = SETBIT{Key: key, Offset: 20, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))

		value, err = GETBIT{Key: key, Offset: 20}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))
	}
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
		value, isNil, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, "")
		assert.Equal(t, isNil, true)
		assert.NoError(t, err)
	}
	{// SET key value
		newValue := "nimo" + strconv.Itoa(time.Now().Nanosecond())
		_, _, err = client.DoStringReply(ctx, []string{"SET", key, newValue}) ; assert.NoError(t, err)
		_="GET key"
		value, isNil, err := GET{
			Key: key,
		}.Do(ctx, client)
		assert.Equal(t, value, newValue)
		assert.Equal(t, isNil, false)
		assert.NoError(t, err)
	}

}


func TestGetBit(t *testing.T) {
	for _, client := range Connecters {
		redisGetBit(t, client)
	}
}
func redisGetBit(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "getbit"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		value, err := GETBIT{Key: key, Offset: 0}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))
	}
	{
		value, err := SETBIT{Key: key, Offset: 0, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))

		value, err = GETBIT{Key: key, Offset: 0}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))
	}
	{
		value, err := GETBIT{Key: key, Offset: 10}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))
	}
	{
		value, err := SETBIT{Key: key, Offset: 10, Value: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(0))

		value, err = GETBIT{Key: key, Offset: 10}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, uint8(1))
	}
}


func TestGetDel(t *testing.T) {
	for _, client := range Connecters {
		redisGetDel(t, client)
	}
}
func redisGetDel(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "getdel"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	isNil, err := SET{Key: key, Value: "Hello",NeverExpire:true}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, isNil, false)
	{
		value,isNil, err := GETDEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "Hello")
	}
	{
		value,isNil, err := GETDEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, value, "")
	}
	{
		_, err := DEL{Key: key+"inexistence"}.Do(ctx, client) ; assert.NoError(t, err)
		value,isNil, err := GETDEL{Key: key+"inexistence"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, value, "")
	}

}

func TestGetEx(t *testing.T) {
	for _, client := range Connecters {
		redisGetEx(t, client)
	}
}
func redisGetEx(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "getex"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		value, isNil, err := GETEX{Key: key, Expire: time.Second}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, value, "")
	}
	{
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		isNil, err := SET{Key: key, Value: "hi",NeverExpire:true}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		value, isNil, err := GETEX{Key: key, Expire: time.Second}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, value, "hi")
		result, err := PTTL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, result.TTL.Milliseconds() > 900, true)
		assert.Equal(t, result.TTL.Milliseconds() <= 1000, true)
	}
}

func TestGetRange(t *testing.T) {
	for _, client := range Connecters {
		redisGetRange(t, client)
	}
}
func redisGetRange(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "getrange"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	_, err = SET{NeverExpire: true, Key: key, Value: "This is a string"}.Do(ctx, client) ; assert.NoError(t, err)
	value, err := GETRANGE{Key: key, Start: 0, End: 3}.Do(ctx, client) ; if err != nil {
	    return
	}
	assert.Equal(t, value, "This")
	{
		value, err := GETRANGE{Key: key, Start: -3, End: -1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, "ing")
	}
	{
		value, err := GETRANGE{Key: key, Start: 0, End: -1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, "This is a string")
	}
	{
		value, err := GETRANGE{Key: key, Start: 0, End: -1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, "This is a string")
	}
	{
		value, err := GETRANGE{Key: key, Start: 10, End: 100}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, value, "string")
	}
}


func TestGetSet(t *testing.T) {
	for _, client := range Connecters {
		redisGetSet(t, client)
	}
}
func redisGetSet(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "getset"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		oldValue, isNil, err := GETSET{Key: key, Value: "a"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, true)
		assert.Equal(t, oldValue, "")
	}
	{
		oldValue, isNil, err := GETSET{Key: key, Value: "b"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, isNil, false)
		assert.Equal(t, oldValue, "a")
	}

}


func TestIncr(t *testing.T) {
	for _, client := range Connecters {
		redisIncr(t, client)
	}
}
func redisIncr(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "incr"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		newValue, err := INCR{Key:key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(1))
	}
	{
		newValue, err := INCR{Key:key}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(2))
	}
}
func TestIncrBy(t *testing.T) {
	for _, client := range Connecters {
		redisIncrBy(t, client)
	}
}
func redisIncrBy(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "incrby"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		newValue, err := INCRBY{Key:key,Increment: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(1))
	}
	{
		newValue, err := INCRBY{Key:key,Increment: 1}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(2))
	}
	{
		newValue, err := INCRBY{Key:key,Increment: 2}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(4))
	}
	{
		newValue, err := INCRBY{Key:key,Increment: -2}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, int64(2))
	}
}

func TestIncrByFloat(t *testing.T) {
	for _, client := range Connecters {
		redisIncrByFloat(t, client)
	}
}
func redisIncrByFloat(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "incrbyfloat"
	_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
	{
		newValue, err := INCRBYFLOAT{Key:key,Increment: "1.1"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(1.1))
	}
	{
		newValue, err := INCRBYFLOAT{Key:key,Increment: "1.3"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(2.4))
	}
	{
		newValue, err := INCRBYFLOAT{Key:key,Increment: "-0.1"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(2.3))
	}
	{
		newValue, err := INCRBYFLOAT{Key:key,Increment: "-.1"}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(2.2))
	}
	{
		value := strconv.FormatFloat(1.333333, 'f', 2, 64)
		newValue, err := INCRBYFLOAT{Key:key,Increment: value}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(3.53))
	}
	{
		value := strconv.FormatFloat(1.333333, 'f', 5, 64)
		newValue, err := INCRBYFLOAT{Key:key,Increment: value}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, newValue, float64(4.86333))
	}
}


func TestMGet(t *testing.T) {
	for _, client := range Connecters {
		redisMGet(t, client)
	}
}
func redisMGet(t *testing.T, client Connecter) {
	ctx := context.TODO()
	keys := []string{"k1", "k2"}
	_, err := DEL{Keys: keys}.Do(ctx, client) ; assert.NoError(t, err)
	values, err := MGET{Keys: keys}.Do(ctx, client) ; assert.NoError(t, err)
	assert.Equal(t, values, ArrayString{
		{Valid: false,String:""},
		{Valid: false,String:""},
	})
	{
		_, err := SET{NeverExpire: true, Key: keys[0], Value:"a"}.Do(ctx, client) ; assert.NoError(t, err)
		values, err := MGET{Keys: keys}.Do(ctx, client) ; assert.NoError(t, err)
		assert.Equal(t, values, ArrayString{
			{Valid: true,String:"a"},
			{Valid: false,String:""},
		})
	}
}