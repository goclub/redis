package red


import "time"


func Duration (duration time.Duration) OptionDuration {
	return OptionDuration{
		duration: duration,
		valid: true,
	}
}

type OptionDuration struct {
	duration time.Duration
	valid bool
}
func (o OptionDuration) Unwrap() time.Duration {
	return o.duration
}


func Uint (i uint) OptionUint {
	return OptionUint{
		uint: i,
		valid: true,
	}
}

type OptionUint struct {
	uint uint
	valid bool
}
func (o OptionUint) Unwrap() uint {
	return o.uint
}


func Uint32 (i uint32) OptionUint32 {
	return OptionUint32{
		uint32: i,
		valid: true,
	}
}
type OptionUint32 struct {
	uint32 uint32
	valid bool
}
func (o OptionUint32) Unwrap() uint32 {
	return o.uint32
}


func Uint8 (i uint8) OptionUint8 {
	return OptionUint8{
		uint8: i,
		valid: true,
	}
}
type OptionUint8 struct {
	uint8 uint8
	valid bool
}
func (o OptionUint8) Unwrap() uint8 {
	return o.uint8
}


func Uint64 (i uint64) OptionUint64 {
	return OptionUint64{
		uint64: i,
		valid: true,
	}
}
type OptionUint64 struct {
	uint64 uint64
	valid bool
}
func (o OptionUint64) Unwrap() uint64 {
	return o.uint64
}