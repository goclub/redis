package red

import (
	"context"
	xerr "github.com/goclub/error"
	xtime "github.com/goclub/time"
	"strconv"
	"time"
)
type COPY struct {
	Source string
	Destination string
	DB OptionUint8
	Replace bool
}
// When Source is exists and Destination is inexistence return 1 Otherwise return 0
func (data COPY) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	args := []string{"COPY"}

	args = append(args, data.Source, data.Destination)
	if data.DB.Valid {
		args = append(args, "DB")
		args = append(args, strconv.FormatUint(uint64(data.DB.Uint8), 10))
	}
	if data.Replace {
		args = append(args, "REPLACE")
	}
	reply, _, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type DUMP struct {
	Key string
}
func (data DUMP) Do(ctx context.Context, client Connecter) (value string, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"DUMP", data.Key}
	value,_, err = client.DoStringReply(ctx, args) ; if err != nil {
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
	if len(data.Keys) == 0 { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args = append(args, data.Keys...)
	var value int64
	value,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	delCount = uint64(value)
	return
}


type EXISTS struct {
	Key string
	Keys []string
}
func (data EXISTS) Do(ctx context.Context, client Connecter) (existsCount uint64, err error) {
	args := []string{"EXISTS"}
	if data.Key != "" {
		data.Keys = []string{data.Key}
	}
	if len(data.Keys) == 0 { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args = append(args, data.Keys...)
	var value int64
	value,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	existsCount = uint64(value)
	return
}

type KEYS struct {
	Pattern string
}
func (data KEYS) Do(ctx context.Context, client Connecter) (keys []string, err error) {
	args := []string{"KEYS", data.Pattern}
	reply, err := client.DoArrayStringReply(ctx, args) ; if err != nil {
		return
	}
	for _, v := range reply {
		keys = append(keys, v.String)
	}
	return
}
type PEXPIRE struct {
	Key string
	Duration time.Duration
	NX bool
	XX bool
	GT bool
	LT bool
}

func (data PEXPIRE) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	if data.Duration < time.Millisecond-1 {
		err = xerr.New("red.PEXPIRE{}.Duration can not less than time.Millisecond")
		return
	}
	args := []string{"PEXPIRE", data.Key, strconv.FormatInt(data.Duration.Milliseconds(), 10)}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	if data.GT {
		args = append(args, "GT")
	}
	if data.LT {
		args = append(args, "LT")
	}
	reply,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type EXPIRE struct {
	Key string
	Duration time.Duration
	NX bool
	XX bool
	GT bool
	LT bool
}
func (data EXPIRE) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	if data.Duration < time.Second-1 {
		err = xerr.New("red.EXPIRE{}.Duration can not less than time.Second")
		return
	}
	args := []string{"EXPIRE", data.Key, strconv.FormatUint(uint64(data.Duration.Seconds()), 10)}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	if data.GT {
		args = append(args, "GT")
	}
	if data.LT {
		args = append(args, "LT")
	}
	reply,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type EXPIREAT struct {
	Key string
	Time time.Time
	NX bool
	XX bool
	GT bool
	LT bool
}
func (data EXPIREAT) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"EXPIREAT", data.Key, strconv.FormatInt(xtime.UnixMilli(data.Time), 10)}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	if data.GT {
		args = append(args, "GT")
	}
	if data.LT {
		args = append(args, "LT")
	}
	reply,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	return
}
type EXPIRETIME struct {
	Key string
}
func (data EXPIRETIME) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
	args := []string{"EXPIRETIME", data.Key}
	reply,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
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
	if data.Key == "" { err = xerr.New("goclub/redis: key can not be empty string") ; return}
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