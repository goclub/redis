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
