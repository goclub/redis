package red

import (
	"strconv"
	"strings"
)

// >= 6.0: timeout is interpreted as a double instead of an integer.
type Second struct {
	Value float64
}
func (data Second) String() string {
	s := strconv.FormatFloat(data.Value, 'f', 2,64)
	return strings.TrimSuffix(s, `.00`)
}