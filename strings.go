package red

import (
	"context"
	"errors"
	xerr "github.com/goclub/error"
	xtime "github.com/goclub/time"
	"strconv"
	"time"
)
type APPEND struct {
	Key string
	Value string
}
func (data APPEND) Do(ctx context.Context, client Connecter) (length uint64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"APPEND", data.Key, data.Value}
	value, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
	    return
	}
	length = uint64(value)
	return
}

type BITCOUNT struct {
	Key string
	// start and end offset unit byte (8bit)
	StartByte OptionInt64
	EndByte OptionInt64
}

func (data BITCOUNT) Do(ctx context.Context, client Connecter) (length uint64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"BITCOUNT", data.Key}
	if data.StartByte.Valid {
		args = append(args, strconv.FormatInt(data.StartByte.Int64, 10))
	}
	if data.EndByte.Valid {
		args = append(args, strconv.FormatInt(data.EndByte.Int64, 10))
	}
	value, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	length = uint64(value)
	return
}

type BITFIELD struct {
	Key string
	Args []string
}

func (data BITFIELD) Do(ctx context.Context, client Connecter) (reply []OptionInt64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"BITFIELD", data.Key}
	args = append(args, data.Args...)
	return client.DoArrayIntegerReply(ctx, args)
}

// bit operation
type BITOP struct {
	AND bool
	OR bool
	XOR bool
	NOT bool
	DestKey string
	Key string
	Keys []string
}
func (data BITOP) Do(ctx context.Context, client Connecter) (size uint64, err error) {
	args := []string{"BITOP"}
	if data.AND {
		args = append(args, "AND")
	}
	if data.OR {
		args = append(args, "OR")
	}
	if data.XOR {
		args = append(args, "XOR")
	}
	if data.NOT {
		args = append(args, "NOT")
	}
	args = append(args, data.DestKey)
	if data.Key != "" {
		data.Keys = []string{data.Key}
	}
	if len(data.Keys) == 0 { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args = append(args, data.Keys...)
	value,_, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	size = uint64(value)
	return
}

type BITPOS struct {
	Key string
	Bit uint8
	Start OptionUint64
	End OptionUint64
}

func (data BITPOS) Do(ctx context.Context, client Connecter) (position int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"BITPOS", data.Key, strconv.FormatUint(uint64(data.Bit), 10)}
	if data.Start.Valid {
		args = append(args, strconv.FormatUint(data.Start.Uint64, 10))
	}
	if data.End.Valid {
		args = append(args, strconv.FormatUint(data.End.Uint64, 10))
	}
	position, _, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type DECR struct {
	Key string
}
func (data DECR) Do(ctx context.Context, client Connecter) (newValue int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"DECR", data.Key}
	newValue,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type DECRBY struct {
	Key string
	Decrement int64
}
func (data DECRBY) Do(ctx context.Context, client Connecter) (newValue int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"DECRBY", data.Key, strconv.FormatInt(data.Decrement, 10)}
	newValue,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}

type GET struct {
	Key string
}
func (data GET) Do(ctx context.Context, client Connecter) (value string, isNil bool, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"GET", data.Key}
	return client.DoStringReply(ctx, args)
}

type GETBIT struct {
	Key string
	Offset uint64
}

func (data GETBIT) Do(ctx context.Context, client Connecter) (value uint8, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"GETBIT", data.Key, strconv.FormatUint(data.Offset, 10)}
	reply, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	value = uint8(reply)
	return
}

type GETDEL struct {
	Key string
}
func (data GETDEL) Do(ctx context.Context, client Connecter) (value string,isNil bool, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"GETDEL", data.Key}
	return client.DoStringReply(ctx, args)
}
type GETEX struct {
	Key string
	Expire time.Duration
	ExpireAt time.Time
	PERSIST bool
}
func (data GETEX) Do(ctx context.Context, client Connecter) (value string,isNil bool, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"GETEX", data.Key}
	if data.Expire != 0 {
		px := strconv.FormatInt(data.Expire.Milliseconds(), 10)
		args = append(args, "PX", px)
	}
	if data.ExpireAt.IsZero() == false {
		args = append(args, "PXAT", strconv.FormatInt(xtime.UnixMilli(data.ExpireAt), 10))
	}
	if data.PERSIST {
		args = append(args, "PERSIST")
	}
	return client.DoStringReply(ctx, args)
}

type GETRANGE struct {
	Key string
	Start int64
	End int64
}
func (data GETRANGE) Do(ctx context.Context, client Connecter) (value string, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"GETRANGE", data.Key, strconv.FormatInt(data.Start, 10), strconv.FormatInt(data.End, 10)}
	value, _, err = client.DoStringReply(ctx, args) ; if err != nil {
		return
	}
	return
}

type INCR struct {
	Key string
}
func (data INCR) Do(ctx context.Context, client Connecter) (newValue int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"INCR", data.Key}
	newValue,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type INCRBY struct {
	Key string
	Increment int64
}
func (data INCRBY) Do(ctx context.Context, client Connecter) (newValue int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"INCRBY", data.Key, strconv.FormatInt(data.Increment, 10)}
	newValue,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}

type INCRBYFLOAT struct {
	Key string
	Increment string `eg:"strconv.FormatFloat(value, 'f', 2, 64)"`
}
func (data INCRBYFLOAT) Do(ctx context.Context, client Connecter) (newValue float64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"INCRBYFLOAT", data.Key, data.Increment}
	reply, _, err := client.DoStringReply(ctx, args) ; if err != nil {
		return
	}
	return strconv.ParseFloat(reply, 64)
}
type MGET struct {
	Keys []string
}
func (data MGET) Do(ctx context.Context, client Connecter) (values []OptionString, err error) {
	if len(data.Keys) == 0 { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"MGET"}
	args = append(args, data.Keys...)
	values, err = client.DoArrayStringReply(ctx, args) ; if err != nil {
		return
	}
	return

}
type MSET struct {
	KeysValues []KeyValue
}
func (data MSET) Do(ctx context.Context, client Connecter) (err error) {
	args := []string{"MSET"}
	if len(data.KeysValues) == 0 { err = xerr.New("goclub/redis: KeysValues can not be empty slice") ; return}
	for _, KeyValue := range data.KeysValues {
		args = append(args, KeyValue.Key, KeyValue.Value)
	}
	_, _, err = client.DoStringReply(ctx, args) ; if err != nil {
	    return
	}
	return
}
type MSETNX struct {
	KeysValues []KeyValue
}
func (data MSETNX) Do(ctx context.Context, client Connecter) (result int8, err error) {
	if len(data.KeysValues) == 0 { err = xerr.New("goclub/redis: KeysValues can not be empty slice") ; return}
	args := []string{"MSETNX"}
	for _, KeyValue := range data.KeysValues {
		args = append(args, KeyValue.Key, KeyValue.Value)
	}
	reply, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	result = int8(reply)
	return
}
type SET struct {
	NeverExpire bool
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time // >= 6.2: Added the GET, EXAT and PXAT option. (ExpireAt)
	KeepTTL bool // >= 6.0: Added the KEEPTTL option.
	XX bool
	NX bool
	GET bool
}
var ErrSetForgetTimeToLive = errors.New("goclub/redis: SET maybe you forget set field Expire or ExpireAt or KeepTTL or NeverExpire")
func (data SET) Do(ctx context.Context, client Connecter) (reply string, isNil bool ,err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"SET", data.Key, data.Value}
	// 只有在明确 NeverExpire 时候才允许 Expire 留空
	if data.NeverExpire == false && data.Expire == 0 && data.ExpireAt.IsZero() && data.KeepTTL == false {
		err = ErrSetForgetTimeToLive
		return
	}
	if data.Expire != 0 {
		px := strconv.FormatInt(data.Expire.Milliseconds(), 10)
		args = append(args, "PX", px)
	}
	if data.ExpireAt.IsZero() == false {
		args = append(args, "PXAT", strconv.FormatInt(xtime.UnixMilli(data.ExpireAt), 10))
	}
	if data.KeepTTL {
		args = append(args, "KEEPTTL")
	}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	if data.GET {
		args = append(args, "GET")
	}
	return client.DoStringReply(ctx, args)
}
type SETBIT struct {
	Key string
	Offset uint64
	Value uint8
}

func (data SETBIT) Do(ctx context.Context, client Connecter) (originalValue uint8, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"SETBIT", data.Key, strconv.FormatUint(data.Offset, 10), strconv.FormatUint(uint64(data.Value), 10)}
	reply, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	originalValue = uint8(reply)
	return
}


type SETRANGE struct {
	Key string
	Offset int64
	Value string
}
func (data SETRANGE) Do(ctx context.Context, client Connecter) (length int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"SETRANGE", data.Key, strconv.FormatInt(data.Offset, 10), data.Value}
	length, _, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type STRLEN struct {
	Key string
}
func (data STRLEN) Do(ctx context.Context, client Connecter) (length int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"STRLEN", data.Key}
	length, _, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}