package red

import (
	"errors"
	"strings"
	"time"
)

func checkDuration(command string, arg string, duration time.Duration) error {
	if duration == 0 {
		return nil
	}
	if duration < time.Millisecond {
		return errors.New("goclub/redis:(ERR_DURATION) " + command + " "+ arg + "  can not set " + duration.String() + ", maybe you forget time.Millisecond or time.time.Second")
	}
	return nil
}
func checkKey(command string, arg string, key string) error {
	if len(key) == 0 {
		return errors.New(strings.Join([]string{"goclub/redis(ERR_EMPTY_KEY)" , command ,  arg ,"key is empty"}, " "))
	}
	return nil
}