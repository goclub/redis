package red

import (
	"context"
	"strconv"
)

type HGET struct {
	Key string
	Field string
}
func (data HGET) Do(ctx context.Context, client Client) (value string, has bool, err error) {
	cmd := "HGET"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key, data.Field}
	result, err := Command(ctx, client, &value, args) ; if err != nil {
		return
	}
	if result.IsNil == false {
		has = true
	}
	return
}
type HMGET struct {
	Key string
	Fields []string
}
func (data HMGET) Do(ctx context.Context, client Client) (values []string, err error) {
	cmd := "HMGET"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key}
	args = append(args, data.Fields...)
	_, err = Command(ctx, client, &values, args) ; if err != nil {
		return
	}
	return
}

type HSET struct {
	Key string
	Field string
	Value string
	// As of Redis 4.0.0, HSET is variadic and allows for multiple field/value pairs.
	Multiple []FieldValue
}
type FieldValue struct {
	Field string
	Value string
}

func (data HSET) Do(ctx context.Context, client Client) (added uint, err error) {
	cmd := "HSET"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key}
	if len(data.Multiple) == 0 {
		args = append(args, data.Field, data.Value)
	} else {
		var sets []string
		for _, item := range data.Multiple {
			sets = append(sets, item.Field, item.Value)
		}
		args = append(args, sets...)
	}
	_, err = Command(ctx, client, &added, args) ; if err != nil {
		return
	}
	return
}
type HINCRBY struct {
	Key string
	Field string
	Increment uint64
}
func (data HINCRBY) Do(ctx context.Context, client Client) (value uint64, err error) {
	cmd := "HINCRBY"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key,data.Field, strconv.FormatUint(data.Increment, 10)}
	_, err = Command(ctx, client, &value, args) ; if err != nil {
		return
	}
	return
}