package red

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
)

type StreamID struct {
	UnixMilli uint64
	Seq uint64
}
var ErrInvalidStreamID = errors.New("goclub/redis: invalid stream id")

func NewStreamID(s string) (streamID StreamID, err error) {
	data := strings.Split(s, "-")
	if len(data) != 2 {
		return streamID, ErrInvalidStreamID
	}
	unixMilli, err := strconv.ParseUint(data[0], 10, 64) ; if err != nil {
		return
	}
	seq, err := strconv.ParseUint(data[1], 10, 64) ; if err != nil {
		return
	}
	streamID.UnixMilli, streamID.Seq = unixMilli, seq
	return
}

func (id StreamID) Before(pivot StreamID) bool {
	if id.UnixMilli != id.UnixMilli {
		return id.UnixMilli < pivot.UnixMilli
	}
	return id.Seq < pivot.Seq
}
func (id StreamID) After(pivot StreamID) bool {
	if id.UnixMilli != id.UnixMilli {
		return id.UnixMilli > pivot.UnixMilli
	}
	return id.Seq > pivot.Seq
}


var maxUint64Len = len(strconv.FormatUint(math.MaxUint64, 10))

func (id StreamID) Bytes() []byte {
	b := make([]byte, 0, maxUint64Len*2+1)
	b = strconv.AppendUint(b, id.UnixMilli, 10)
	b = append(b, '-')
	b = strconv.AppendUint(b, id.Seq, 10)
	return b
}
func (id StreamID) String() string {
	return string(id.Bytes())
}

// Available since 5.0.0
type XADD struct {
	Key string
	FieldValues []FieldValue
	// id 可留空，留空则为 *
	ID string
}
func (data XADD) Do(ctx context.Context, client Client) (streamID StreamID, err error) {
	cmd := "XADD"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	if len(data.ID) == 0 {
		data.ID = "*"
	}
	args := []string{cmd, data.Key, data.ID,}
	if len(data.FieldValues) == 0 {
		return streamID, errors.New(`goclub/redis: XADD{Key: "` + data.Key + `"} FieldsValue can not empty`)
	}
	for _, item := range data.FieldValues {
		args = append(args, item.Field, item.Value)
	}
	var s string
	_, err = Command(ctx, client, &s, args) ; if err != nil {
		return
	}
	streamID, err = NewStreamID(s) ; if err != nil {
		return
	}
	return
}
// Available since 5.0.0.
type XLEN struct {
	Key string
}
func (data XLEN) Do(ctx context.Context, client Client) (length uint, err error) {
	cmd := "XLEN"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key}
	_, err = Command(ctx, client, &length, args) ; if err != nil {
		return
	}
	return
}
type XRANGE struct {
	Key string
	Start string
	End string
	Count uint64
}
func (data XRANGE) Do(ctx context.Context, client Client, streamEntryEntrySlicePtr interface{}) (err error) {
	cmd := "XRANGE"
	err = checkKey(cmd, "Key", data.Key) ; if err != nil {
		return
	}
	err = checkKey(cmd, "Start", data.Start) ; if err != nil {
		return
	}
	err = checkKey(cmd, "End", data.End) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key, data.Start, data.End}
	if data.Count != 0 {
		args = append(args, "COUNT", strconv.FormatUint(data.Count, 10))
	}
	_, err = Command(ctx, client, streamEntryEntrySlicePtr, args) ; if err != nil {
		return
	}
	return
}
type XDEL struct {
	Key string
	StreamID []string
}
// 通常，你可能将Redis流想象为一个仅附加的数据结构，但是Redis流是存在于内存中的， 所以我们也可以删除条目。这也许会有用，例如，为了遵守特定的隐私策略
func (data XDEL) Do(ctx context.Context, client Client) (delCount uint, err error) {
	cmd := "XDEL"
	err = checkKey(cmd, "Key", data.Key) ; if err != nil {
		return
	}
	if len(data.StreamID) == 0 {
			return 0, errors.New("goclub/redis: XDEL id cannot be empty")
	}
	args := []string{cmd, data.Key}
	args = append(args, data.StreamID...)
	_, err = Command(ctx, client, &delCount, args) ; if err != nil {
		return
	}
	return
}