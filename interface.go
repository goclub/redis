package red

import (
	"context"
)

// Connecter RESP
type Connecter interface {
	DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error)
	DoStringReplyWithoutNil(ctx context.Context, args []string) (reply string, err error)
	DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error)
	DoIntegerReplyWithoutNil(ctx context.Context, args []string) (reply int64, err error)
	DoArrayIntegerReply(ctx context.Context, args []string) (reply []OptionInt64, err error)
	DoArrayStringReply(ctx context.Context, args []string) (reply []OptionString, err error)
	Eval(ctx context.Context, script Script) (reply Reply, isNil bool, err error)
	EvalWithoutNil(ctx context.Context, script Script) (reply Reply, err error)
}

type API struct {
	// strings
	Append      APPEND
	BitCount    BITCOUNT
	BitField    BITFIELD
	BitOp       BITOP
	BitPos      BITPOS
	Decr        DECR
	DecrBy      DECRBY
	Get         GET
	GetBit      GETBIT
	GetDel      GETDEL
	GetEx       GETEX
	GetRange    GETRANGE
	// GETSET: Please use SET (GETSET有清空TTL的"隐患")
	Incr        INCR
	IncrBy      INCRBY
	IncrByFloat INCRBYFLOAT
	MGet        MGET
	MSet        MSET
	MSetNX      MSETNX
	// PSETEX: Please use SET
	Set      SET
	SetBit   SETBIT
	SetRange SETRANGE
	// STRALGO TODO
	StrLen STRLEN

	// keys
	Copy       COPY
	Del        DEL
	Dump       DUMP
	Exists     EXISTS
	Expire     EXPIRE
	ExpireAt   EXPIREAT   // @needtest
	ExpireTime EXPIRETIME // @needtest
	Keys       KEYS
	PExpire    PEXPIRE
	// PEXPIREAT
	// PEXPIRETIME
	PTTL PTTL
	// RENAME
	// RENAMENX
	// RESTORE
	// SORT
	// SORT_RO
	// TOUCH
	// TTL
	// TYPE
	// UNLINK
	// WAIT
	// SCAN
	HDEL HDEL
}
