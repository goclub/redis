package red

import (
	"context"
	"errors"
	xtime "github.com/goclub/time"
	"strconv"
	"time"
)
type APPEND struct {
	Key string
	Value string
}
func (data APPEND) Do(ctx context.Context, client Connecter) (length uint64, err error) {
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
	args := []string{"BITCOUNT", data.Key}
	if data.StartByte.valid {
		args = append(args, strconv.FormatInt(data.StartByte.int64, 10))
	}
	if data.EndByte.valid {
		args = append(args, strconv.FormatInt(data.EndByte.int64, 10))
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

func (data BITFIELD) Do(ctx context.Context, client Connecter) (reply []int64, err error) {
	args := []string{"BITFIELD", data.Key}
	args = append(args, data.Args...)
	reply, _, err = client.DoIntegerSliceReply(ctx, args) ; if err != nil {
		return
	}
	return
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
	args := []string{"BITPOS", data.Key, strconv.FormatUint(uint64(data.Bit), 10)}
	if data.Start.valid {
		args = append(args, strconv.FormatUint(data.Start.uint64, 10))
	}
	if data.End.valid {
		args = append(args, strconv.FormatUint(data.End.uint64, 10))
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
	args := []string{"DECRBY", data.Key, strconv.FormatInt(data.Decrement, 10)}
	newValue,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type DEL struct {
	Key string
	Keys []string
}
func (data DEL) Do(ctx context.Context, client Connecter) (delCount uint64, err error) {
	args := []string{"DEL"}
	if data.Key != "" {
		data.Keys = []string{data.Key}
	}
	args = append(args, data.Keys...)
	var value int64
	value,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	delCount = uint64(value)
	return
}

type GET struct {
	Key string
}

func (data GET) Do(ctx context.Context, client Connecter) (value string, hasValue bool, err error) {
	args := []string{"GET", data.Key}
	value, isNil, err := client.DoStringReply(ctx, args) ; if err != nil {
		// error
		return "", false, err
	}
	if isNil {
		// key 不存在
		return "", false, nil
	} else {
		// key 存在
		return value, true, nil
	}
}



type GETBIT struct {
	Key string
	Offset uint64
}

func (data GETBIT) Do(ctx context.Context, client Connecter) (value uint8, err error) {
	args := []string{"GETBIT", data.Key, strconv.FormatUint(data.Offset, 10)}
	reply, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	value = uint8(reply)
	return
}



type SET struct {
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time // >= 6.2: Added the GET, EXAT and PXAT option. (ExpireAt)
	KeepTTL bool // >= 6.0: Added the KEEPTTL option.
	NeverExpire bool
	XX bool
	NX bool
}
var ErrSetForgetTimeToLive = errors.New("goclub/redis: SET maybe you forget set field Expire or ExpireAt or KeepTTL or NeverExpire")
func (data SET) Do(ctx context.Context, client Connecter) (isNil bool ,err error) {
	args := []string{"SET", data.Key, data.Value}
	// 只有在明确 NeverExpire 时候才允许 Expire 留空
	if data.NeverExpire == false && data.Expire == 0 && data.ExpireAt.IsZero() && data.KeepTTL == false {
		return false, ErrSetForgetTimeToLive
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
	_, isNil, err = client.DoStringReply(ctx, args) ; if err != nil {
		return
	}
	return
}

type SETBIT struct {
	Key string
	Offset uint64
	Value uint8
}

func (data SETBIT) Do(ctx context.Context, client Connecter) (originalValue uint8, err error) {
	args := []string{"SETBIT", data.Key, strconv.FormatUint(data.Offset, 10), strconv.FormatUint(uint64(data.Value), 10)}
	reply, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	originalValue = uint8(reply)
	return
}


type PTTL struct {
	Key string
}
type ResultTTL struct {
	TTL time.Duration
	NeverExpire bool
	KeyDoesNotExist bool
}
func (data PTTL) Do(ctx context.Context, client Connecter) (result ResultTTL, err error) {
	args := []string{"PTTL", data.Key}
	value, _, err := client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	if value == -1 {
		result.NeverExpire = true
		return
	}
	if value == -2 {
		result.KeyDoesNotExist = true
		return
	}
	result.TTL = time.Millisecond * time.Duration(value)
	return
}