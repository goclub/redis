package red

import (
	"context"
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
}

func (data PEXPIRE) Do(ctx context.Context, client Connecter) (reply int64, err error) {
	args := []string{"PEXPIRE", data.Key, strconv.FormatInt(data.Duration.Milliseconds(), 10)}
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