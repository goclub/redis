package red

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

type StreamID struct {
	UnixMilli uint64
	Seq uint64
}
var ErrInvalidStreamID = errors.New("goclub/redis:  invalid stream id")

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
type QueryStream struct {
	Key string
	ID string
}
type XREAD struct {
	Streams []QueryStream
	Block time.Duration
	Count uint
}
func (data XREAD) Do(ctx context.Context, client Client, streamEntrySlicePtr interface{}) (err error) {
	cmd := "XREAD"
	if len(data.Streams) == 0 {
		return errors.New("goclub/redis:  XREAD Streams can not be empty")
	}
	args := []string{cmd}
	if data.Count != 0 {
		args = append(args, "COUNT", strconv.FormatUint(uint64(data.Count), 10))
	}
	if data.Block != 0 {
		args = append(args, "BLOCK", strconv.FormatInt(data.Block.Milliseconds(), 10))
	}
	// STREAMS 必须是最后一个选项
	args = append(args, "STREAMS")
	var streamKeyList []string
	var idList []string
	for _, keyID := range data.Streams {
		streamKeyList = append(streamKeyList, keyID.Key)
		idList = append(idList, keyID.ID)
	}
	args = append(args, streamKeyList...)
	args = append(args, idList...)
	_, err = Command(ctx, client, streamEntrySlicePtr, args) ; if err != nil {

		return
	}
	return
}
// XTRIM key MAXLEN|MINID [=|~] threshold [LIMIT count]
type XTRIM struct {
	Key string
	MaxLen OptionUint
	Tilde bool
	// >= 6.2: Added the MINID trimming strategy and the LIMIT option.
	MinID string
	// >= 6.2: Added the MINID trimming strategy and the LIMIT option.
	LIMIT OptionUint
}
func (data XTRIM) Do(ctx context.Context, client Client) (delCount uint, err error) {
	cmd := "XTRIM"
	err = checkKey(cmd, "", data.Key) ; if err != nil {
		return
	}
	args := []string{cmd, data.Key}
	if data.MaxLen.valid {
		args = append(args, "MAXLEN")
		if data.Tilde {
			args = append(args, "~")
		}
		args = append(args, strconv.FormatUint(uint64(data.MaxLen.uint), 10))
		if data.LIMIT.valid {
			args = append(args, "LIMIT", strconv.FormatUint(uint64(data.LIMIT.uint), 10))
		}
	} else if data.MinID != "" {
		args = append(args, "MINID", data.MinID)
	}
	_, err = Command(ctx, client, &delCount, args) ; if err != nil {
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
func (data XRANGE) Do(ctx context.Context, client Client, streamEntrySlicePtr interface{}) (err error) {
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
	_, err = Command(ctx, client, streamEntrySlicePtr, args) ; if err != nil {
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
			return 0, errors.New("goclub/redis:  XDEL id cannot be empty")
	}
	args := []string{cmd, data.Key}
	args = append(args, data.StreamID...)
	_, err = Command(ctx, client, &delCount, args) ; if err != nil {
		return
	}
	return
}

type XGROUPCreate struct {
	Key string
	Group string
	ID string
	MKSTREAM bool
}
func (data XGROUPCreate) Do(ctx context.Context, client Client) (err error) {
	cmd := "XGROUP"
	err = checkKey(cmd, "CREATE key", data.Key) ; if err != nil {
		return
	}
	err = checkKey(cmd, "CREATE group name", data.Group) ; if err != nil {
		return
	}
	err = checkKey(cmd, "CREATE id", data.ID) ; if err != nil {
		return
	}
	args := []string{cmd, "CREATE", data.Key, data.Group, data.ID}
	if data.MKSTREAM {
		args = append(args, "MKSTREAM")
	}
	_, err = Command(ctx, client, nil, args) ; if err != nil {
		return
	}
	return
}
type StartEndCount struct {
	Start string
	End string
	Count uint64
}
type XPENDING struct {
	Key string
	Group string
	IDLE OptionDuration
	StartEndCount StartEndCount
	Consumer string
}

func (data XPENDING) Do(ctx context.Context, client Client, streamEntrySlicePtr interface{}) (err error) {
	cmd := "XPENDING"
	err = checkKey(cmd, "Key", data.Key) ; if err != nil {
		return
	}
	err = checkKey(cmd, "group", data.Group) ; if err != nil {
		return
	}

	args := []string{cmd, data.Key, data.Group}
	if data.IDLE.valid {
		args = append(args, "IDLE", strconv.FormatInt(data.IDLE.Unwrap().Milliseconds(), 10))
	}
	if data.StartEndCount.Start != "" {
		countString := strconv.FormatUint(data.StartEndCount.Count, 10)
		args = append(args, data.StartEndCount.Start, data.StartEndCount.End, countString)
	}
	if data.Consumer != "" {
		args = append(args, data.Consumer)
	}
	_, err = Command(ctx, client, streamEntrySlicePtr, args) ; if err != nil {
		return
	}
	return
}

type XACK struct {
	Key string
	Group string
	StreamID string
	StreamIDs []string
}
func (data XACK) Do(ctx context.Context, client Client) (ackCount uint64, err error) {
	cmd := "XACK"
	err = checkKey(cmd, "Key", data.Key) ; if err != nil {
		return
	}
	err = checkKey(cmd, "group", data.Group) ; if err != nil {
		return
	}
	if data.StreamID != "" {
		data.StreamIDs = append(data.StreamIDs, data.StreamID)
	}
	if len(data.StreamIDs) == 0 {
		return 0, errors.New("goclub/redis:  red.XACK{} StreamID or StreamIDs cannot be empty")
	}
	args := []string{cmd, data.Key, data.Group}
	args = append(args, data.StreamIDs...)
	_, err = Command(ctx, client, &ackCount, args) ; if err != nil {
		return
	}
	return
}

type XREADGROUP struct {
	Group string
	Consumer string
	Count uint64
	Block time.Duration
	Streams []QueryStream
}

func (data XREADGROUP) Do(ctx context.Context, client Client, streamEntrySlicePtr interface{}) (err error) {
	cmd := "XREADGROUP"
	err = checkKey(cmd, "Key", data.Group) ; if err != nil {
		return
	}
	err = checkKey(cmd, "Group", data.Group) ; if err != nil {
		return
	}
	err = checkKey(cmd, "Consumer", data.Consumer) ; if err != nil {
		return
	}
	if len(data.Streams) == 0 {
		return errors.New("goclub/redis:  red.XREADGROUP{} Streams can not be empty slice")
	}
	args := []string{cmd, "GROUP", data.Group, data.Consumer}
	if data.Count != 0 {
		args = append(args, "COUNT", strconv.FormatUint(data.Count, 10))
	}
	if data.Block != 0 {
		args = append(args, "BLOCK", strconv.FormatInt(data.Block.Milliseconds(), 10))
	}
	args = append(args, "STREAMS")
	var keys []string
	var ids []string
	for _, item := range data.Streams {
		keys = append(keys, item.Key)
		ids = append(ids, item.ID)
	}
	args = append(args, keys...)
	args = append(args, ids...)
	_, err = Command(ctx, client, streamEntrySlicePtr, args) ; if err != nil {
		return
	}
	return
}
