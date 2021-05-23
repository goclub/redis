package red

import (
	"context"
)

// RESP Arrays
type ArrayInteger []OptionInt64
type ArrayString []OptionString
type Connecter interface {
	DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error)
	DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error)
	DoArrayIntegerReply(ctx context.Context, args []string)(reply ArrayInteger, err error)
	DoArrayStringReply(ctx context.Context, args []string)(reply ArrayString, err error)
	Eval(ctx context.Context, script Script) (reply interface{}, isNil bool, err error)
}


type API struct {
	Append APPEND
	BitCount BITCOUNT
	BitField BITFIELD
	BitOp BITOP
	BitPos BITPOS
	Decr DECR
	DecrBy DECRBY
	Get GET
	GetBit GETBIT
	GetDel GETDEL
	GetEx GETEX
	GetRange GETRANGE
	GetSet GETSET
	Incr INCR
	IncrBy INCRBY
	IncrByFloat INCRBYFLOAT
	MGet MGET
	MSet MSET
	MSetNX MSETNX
	// PSETEX: Please use SET
	Set SET
	SetBit SETBIT
	SetRange SETRANGE
	// STRALGO TODO
	StrLen STRLEN
	Del DEL
	PTTL PTTL
}