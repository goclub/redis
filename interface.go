package red

import "context"

type Connecter interface {
	DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error)
	DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error)
	DoIntegerSliceReply(ctx context.Context, args []string)(reply []int64, isNil bool, err error)
	Eval(ctx context.Context, script Script) (reply interface{}, isNil bool, err error)
}

type API struct {
	APPEND APPEND
	BITCOUNT BITCOUNT
	BITFIELD BITFIELD
	BITOP BITOP
	SET SET
	GET GET
	DEL DEL
	PTTL PTTL
}