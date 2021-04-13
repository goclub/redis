package red

import (
	"context"
	"errors"
	xtime "github.com/goclub/time"
	"strconv"
	"time"
)


type DEL struct {
	Key string
	Keys []string
}
func (data DEL) Do(ctx context.Context, client Connecter) (delCount uint, err error) {
	args := []string{"DEL"}
	if data.Key != "" {
		data.Keys = append(data.Keys, data.Key)
	}
	args = append(args, data.Keys...)
	var value int64
	value,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	delCount = uint(value)
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

type SET struct {
	Key string
	Value string
	Expire time.Duration
	ExpireAt time.Time // >= 6.2: Added the GET, EXAT and PXAT option. (ExpireAt)
	KeepTTL bool // >= 6.0: Added the KEEPTTL option.
	NeverExpire bool
	NX bool
	XX bool
}
var ErrSetForgetTimeToLive = errors.New("goclub/redis: SET maybe you forget set field Expire or ExpireAt or KeepTTL")
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
	_,isNil, err = client.DoStringReply(ctx, args) ; if err != nil {
		return
	}
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