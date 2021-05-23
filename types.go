package red

import (
	"time"
)


func NewOptionDuration (duration time.Duration) OptionDuration {
	return OptionDuration{
		Valid: true,
		Duration: duration,
	}
}

type OptionDuration struct {
	Valid bool
	Duration time.Duration
}




func NewOptionInt64 (i int64) OptionInt64 {
	return OptionInt64{
		Valid: true,
		Int64: i,
	}
}

type OptionInt64 struct {
	Valid bool
	Int64 int64
}


func NewOptionUint8(i uint8) OptionUint8 {
	return OptionUint8 {
		Valid: true,
		Uint8: i,
	}
}
type OptionUint8 struct {
	Valid bool
	Uint8 uint8
}

func NewOptionUint32 (i uint32) OptionUint32 {
	return OptionUint32{
		Valid: true,
		Uint32: i,
	}
}
type OptionUint32 struct {
	Valid bool
	Uint32 uint32
}


func NewOptionUint64 (i uint64) OptionUint64 {
	return OptionUint64{
		Valid: true,
		Uint64: i,
	}
}
type OptionUint64 struct {
	Valid bool
	Uint64 uint64
}
func NewOptionString(s string) OptionString {
	return OptionString{
		Valid: true,
		String: s,
	}
}
type OptionString struct {
	Valid bool
	String string
}
