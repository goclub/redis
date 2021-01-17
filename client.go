package red

import (
	"context"
	radix4 "github.com/mediocregopher/radix/v4"
	"strconv"
	"time"
)
type Radix4Client struct {
	radix4.Client
}
func NewRadix4Client(radixClient radix4.Client) *Radix4Client {
	return &Radix4Client{radixClient}
}
func (c Radix4Client) RedisExec(ctx context.Context, valuePtr interface{}, cmd string, args []string) (result Result, err error){
	data := radix4.Maybe{
		Rcv: valuePtr,
	}
	err = c.Do(ctx, radix4.Cmd(&data, cmd, args...)) ; if err != nil {
		return
	}
	result = Result{
		IsNil: data.Null,
		IsEmpty: data.Empty,
	}
	return
}
type Execer interface {
	RedisExec (ctx context.Context, valuePtr interface{}, cmd string, args []string) (result Result, err error)
}
type Result struct {
	IsNil bool
	IsEmpty bool
}
type SET struct {
	Key string
	Value string
	Expires time.Duration
	NX bool
	XX bool
}
func (data SET) Do(ctx context.Context, execer Execer) (result Result,err error) {
	args := []string{data.Key, data.Value}
	if data.Expires != 0 {
		px := strconv.FormatInt(int64(data.Expires / time.Millisecond), 10)
		args = append(args, "PX", px)
	}
	if data.NX {
		args = append(args, "NX")
	}
	if data.XX {
		args = append(args, "XX")
	}
	return execer.RedisExec(ctx, nil, "SET", args)
}
type GET struct {
	Key string
}
func (data GET) Do(ctx context.Context, execer Execer) (value string, hasValue bool ,err error) {
	args := []string{data.Key}
	result, err := execer.RedisExec(ctx, &value, "GET", args) ; if err != nil {
		return "", false, err
	}
	if result.IsNil {
		return "", false, err
	} else {
		return value, true, nil
	}
}