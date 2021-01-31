package red

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// >= 6.0: timeout is interpreted as a double instead of an integer.
type Second struct {
	Value float64
}
func (data Second) String() string {
	s := strconv.FormatFloat(data.Value, 'f', 2,64)
	return strings.TrimSuffix(s, `.00`)
}

type Retry struct {
	Times uint8
	Duration time.Duration
}

func (data Retry) check() error {
	if data.Times > 0 && data.Duration == 0 {
		return errors.New("goclub/redis(ERR_MISSING_RETRY_DURATION) if Retry{}.Times > 0 then Retry{}.Duration can not be zero")
	}
	return nil
}