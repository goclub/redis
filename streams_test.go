package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestXADD_XLEN_XREAD_Do(t *testing.T) {
	ctx := context.Background()
	key := "test_stream"
	var err error
	_, err  = red.DEL{Key:key}.Do(ctx, radixClient);
	assert.NoError(t, err)
	_=key
	// 请求前的参数检查
	{
		_ ,err := red.XADD{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_FORGET_ARGS) XADD Key can not be empty")
	}
	{
		_ ,err := red.XADD{Key: key}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, `goclub/redis: XADD{Key: "test_stream"} FieldsValue can not empty`)
	}
	{
		_ ,err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{
				{"name", ""},
			},
		}.Do(ctx, Test{t, "XADD test_stream * name "})
		assert.EqualError(t, err, "goclub/redis: invalid stream id")
	}
	{
		_ ,err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{
				{"name", "nimoc"},
			},
		}.Do(ctx, Test{t, "XADD test_stream * name nimoc"})
		assert.EqualError(t, err, "goclub/redis: invalid stream id")
	}
	{
		_ ,err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{
				{"name", "nimoc"},
				{"age", "18"},
			},
		}.Do(ctx, Test{t, "XADD test_stream * name nimoc age 18"})
		assert.EqualError(t, err, "goclub/redis: invalid stream id")
	}
	{
		_ ,err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{
				{"name", "nimoc"},
				{"age", "18"},
			},
		}.Do(ctx, Test{t, "XADD test_stream * name nimoc age 18"})
		assert.EqualError(t, err, "goclub/redis: invalid stream id")
	}
	{
		_ ,err := red.XLEN{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_FORGET_ARGS) XLEN Key can not be empty")
	}
	{
		err := red.XRANGE{}.Do(ctx, Test{t, ""}, nil)
		assert.EqualError(t, err, `goclub/redis(ERR_FORGET_ARGS) XRANGE Key can not be empty`)
	}
	{
		err := red.XRANGE{Key: key}.Do(ctx, Test{t, ""}, nil)
		assert.EqualError(t, err, `goclub/redis(ERR_FORGET_ARGS) XRANGE Start can not be empty`)
	}
	{
		err := red.XRANGE{Key: key, Start:"-"}.Do(ctx, Test{t, ""}, nil)
		assert.EqualError(t, err, `goclub/redis(ERR_FORGET_ARGS) XRANGE End can not be empty`)
	}
	{
		err := red.XRANGE{Key: key, Start:"-", End: "+"}.Do(ctx, Test{t, "XRANGE test_stream - +"}, nil)
		assert.NoError(t, err)
	}
	{
		err := red.XRANGE{Key: key, Start:"-", End: "+", Count: 1}.Do(ctx, Test{t, "XRANGE test_stream - + COUNT 1"}, nil)
		assert.NoError(t, err)
	}
	{
		_, err := red.XDEL{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis(ERR_FORGET_ARGS) XDEL Key can not be empty")
	}
	{
		_, err := red.XDEL{Key: key}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err, "goclub/redis: XDEL id cannot be empty")
	}
	{
		_, err := red.XDEL{Key:key , StreamID: []string{"1538561700640-0"}}.Do(ctx, Test{t, "XDEL test_stream 1538561700640-0"})
		assert.NoError(t, err)
	}
	{
		_, err := red.XDEL{Key:key , StreamID: []string{"1538561700640-0", "1538561700640-1"}}.Do(ctx, Test{t, "XDEL test_stream 1538561700640-0 1538561700640-1"})
		assert.NoError(t, err)
	}
	{
		err := red.XREAD{}.Do(ctx, Test{t, ""}, nil)
		assert.EqualError(t, err, "goclub/redis: XREAD Streams can not be empty")
	}
	{
		err := red.XREAD{
			Streams: []red.QueryStream{
				{"stream_key_1", "1538561700640-0"},
				{"stream_key_2", "1538561701641-1"},
			},
		}.Do(ctx, Test{t, "XREAD STREAMS stream_key_1 stream_key_2 1538561700640-0 1538561701641-1"}, nil)
		assert.NoError(t, err)
	}
	{
		err := red.XREAD{
			Count: 1,
			Streams: []red.QueryStream{
				{"stream_key_1", "1538561700640-0"},
				{"stream_key_2", "1538561701641-1"},
			},
		}.Do(ctx, Test{t, "XREAD COUNT 1 STREAMS stream_key_1 stream_key_2 1538561700640-0 1538561701641-1"}, nil)
		assert.NoError(t, err)
	}
	{
		err := red.XREAD{
			Block: time.Second,
			Streams: []red.QueryStream{
				{"stream_key_1", "1538561700640-0"},
				{"stream_key_2", "1538561701641-1"},
			},
		}.Do(ctx, Test{t, "XREAD BLOCK 1000 STREAMS stream_key_1 stream_key_2 1538561700640-0 1538561701641-1"}, nil)
		assert.NoError(t, err)
	}
	{
		err := red.XREAD{
			Count: 2,
			Block: time.Second,
			Streams: []red.QueryStream{
				{"stream_key_1", "1538561700640-0"},
				{"stream_key_2", "1538561701641-1"},
			},
		}.Do(ctx, Test{t, "XREAD COUNT 2 BLOCK 1000 STREAMS stream_key_1 stream_key_2 1538561700640-0 1538561701641-1"}, nil)
		assert.NoError(t, err)
	}
	{
		_, err := red.XTRIM{}.Do(ctx, Test{t, ""})
		assert.EqualError(t, err , "goclub/redis(ERR_FORGET_ARGS) XTRIM Key can not be empty")
	}
	{
		_, err := red.XTRIM{
			Key: key,
		}.Do(ctx, Test{t, "XTRIM test_stream"})
		assert.NoError(t ,err)
	}
	{
		_, err := red.XTRIM{
			Key: key,
			MaxLen: red.Uint(100),
		}.Do(ctx, Test{t, "XTRIM test_stream MAXLEN 100"})
		assert.NoError(t ,err)
	}
	{
		_, err := red.XTRIM{
			Key:    key,
			MaxLen: red.Uint(100),
			Tilde:  true,
		}.Do(ctx, Test{t, "XTRIM test_stream MAXLEN ~ 100"})
		assert.NoError(t, err)
	}
	{
		_, err := red.XTRIM{
			Key: key,
			MaxLen: red.Uint(100),
			Tilde: true,
			LIMIT: red.Uint(10),
		}.Do(ctx, Test{t, "XTRIM test_stream MAXLEN ~ 100 LIMIT 10"})
		assert.NoError(t ,err)
	}
	{
		_, err := red.XTRIM{
			Key: key,
			MinID: "1538561700640-0",
		}.Do(ctx, Test{t, "XTRIM test_stream MINID 1538561700640-0"})
		assert.NoError(t ,err)
	}
	{
		// 准备数据
		_, err := red.DEL{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
	}
	{
		// XLEN
		length, err := red.XLEN{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, length, uint(0))
	}
	{
		// XADD
		streamID, err := red.XADD{
			Key: key, 
			FieldValues: []red.FieldValue{
				{"name", "nimoc"},
			},
		}.Do(ctx, radixClient)
		testStreamID(t, streamID)
		assert.NoError(t, err)
	}
	{
		// XLEN
		length, err := red.XLEN{Key: key}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, length, uint(1))
	}
	{
		// XRANGE
		data := []radix.StreamEntry{}
		err = red.XRANGE{
			Key: key,
			Start: "-",
			End: "+",
		}.Do(ctx, radixClient, &data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 1)
		testRadixStreamEntry(t, data[0], [][2]string{{"name","nimoc"}})
	}
	{
		// XADD
		streamID, err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{
				{"name", "goclub"},
			},
		}.Do(ctx, radixClient)
		testStreamID(t, streamID)
		assert.NoError(t, err)
		// XRANGE
		data := []radix.StreamEntry{}
		err = red.XRANGE{
			Key: key,
			Start: "-",
			End: "+",
		}.Do(ctx, radixClient, &data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 2)
		testRadixStreamEntry(t, data[1], [][2]string{{"name","nimoc"}})
		testRadixStreamEntry(t, data[1], [][2]string{{"name","goclub"}})
	}
	{
		// XRANGE - + COUNT 1
		data := []radix.StreamEntry{}
		err = red.XRANGE{
			Key: key,
			Start: "-",
			End: "+",
			Count: 1,
		}.Do(ctx, radixClient, &data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 1)
		testRadixStreamEntry(t, data[0], [][2]string{{"name","nimoc"}})
	}
	{
		// XADD
		streamID, err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{{"name", "nimoc"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		testStreamID(t, streamID)
		// XDEL
		delCount, err := red.XDEL{
			Key:key,
			StreamID: []string{streamID.String()},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(1))
		// XDEL
		delCount, err = red.XDEL{
			Key:key,
			StreamID: []string{streamID.String()},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(0))
		// XADD
		a, err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{{"name", "nimoc"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		testStreamID(t, streamID)
		// XADD
		b, err := red.XADD{
			Key: key,
			FieldValues: []red.FieldValue{{"name", "nimoc"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		testStreamID(t, streamID)
		delCount, err = red.XDEL{
			Key:key,
			StreamID: []string{a.String(), b.String()},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint(2))
	}
	{
		key1 := "test_xread_1"
		key2 := "test_xread_2"
		_, err  = red.DEL{Keys:[]string{key1, key2}}.Do(ctx, radixClient);
		assert.NoError(t, err)
		_, err := red.XADD{
			Key: key1,
			FieldValues: []red.FieldValue{{"name","test_xtrim"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		// XTRIM
		_, err = red.XTRIM{
			Key: key1,
			MaxLen: red.Uint(0),
		}.Do(ctx, radixClient)
		assert.NoError(t ,err)
		_, err = red.XTRIM{
			Key: key2,
			MaxLen: red.Uint(0),
		}.Do(ctx, radixClient)
		assert.NoError(t ,err)
		// XADD key1 XADD key2
		_, err = red.XADD{
			Key:key1,
			FieldValues: []red.FieldValue{{"name","1a"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		_, err = red.XADD{
			Key:key2,
			FieldValues: []red.FieldValue{{"name","2a"}},
		}.Do(ctx, radixClient)
		assert.NoError(t, err)
		{

			// XREAD key1
			var result []radix.StreamEntries
			err = red.XREAD{
				Streams: []red.QueryStream{
					{Key: key1, ID: "0-0"},
				},
			}.Do(ctx, radixClient,  &result)
			assert.NoError(t, err)
			assert.Equal(t, len(result), 1)
			assert.Equal(t, result[0].Stream, "test_xread_1")
			assert.Equal(t, len(result[0].Entries), 1)
			assert.Greater(t, result[0].Entries[0].ID.Time, uint64(1613240335325))
			assert.Equal(t, result[0].Entries[0].ID.Seq, uint64(0))
			assert.Equal(t, result[0].Entries[0].Fields, [][2]string([][2]string{[2]string{"name", "1a"}}))
		}
		{
			// XREAD key1 key2
			var result []radix.StreamEntries
			err = red.XREAD{
				Streams: []red.QueryStream{
					{Key: key1, ID: "0-0"},
					{Key: key2, ID: "0-0"},
				},
			}.Do(ctx, radixClient,  &result)
			assert.NoError(t, err)
			assert.Equal(t, len(result), 2)
			assert.Equal(t, result[0].Stream, "test_xread_1")
			assert.Equal(t, len(result[0].Entries), 1)
			assert.Greater(t, result[0].Entries[0].ID.Time, uint64(1613240335325))
			assert.Equal(t, result[0].Entries[0].ID.Seq, uint64(0))
			assert.Equal(t, result[0].Entries[0].Fields, [][2]string([][2]string{[2]string{"name", "1a"}}))

			assert.Equal(t, result[1].Stream, "test_xread_2")
			assert.Equal(t, len(result[1].Entries), 1)
			assert.Greater(t, result[1].Entries[0].ID.Time, uint64(1613240335325))
			assert.Equal(t, result[1].Entries[0].ID.Seq, uint64(0))
			assert.Equal(t, result[1].Entries[0].Fields, [][2]string([][2]string{[2]string{"name", "2a"}}))
		}
	}
}
func testRadixStreamEntry (t *testing.T, item radix.StreamEntry, fields... [][2]string) {
	assert.Regexp(t, `^\d+$`, item.ID.Time)
	assert.Regexp(t, `^\d$`, item.ID.Seq)
	assert.Equal(t, item.Fields, item.Fields)
}
func testStreamID(t *testing.T, streamID red.StreamID) {
	assert.Regexp(t, `^\d+$`, streamID.UnixMilli)
	assert.Regexp(t, `^\d$`, streamID.Seq)
}

func TestStreamID(t *testing.T) {
	{
		_, err := red.NewStreamID("")
		assert.EqualError(t, err, "goclub/redis: invalid stream id")
	}
	{
		id, err := red.NewStreamID("1526985054069-0")
		assert.NoError(t, err)
		assert.Equal(t, id.UnixMilli, uint64(1526985054069))
		assert.Equal(t, id.Seq, uint64(0))
	}
	{
		id, err := red.NewStreamID("1526985054069-1")
		assert.NoError(t, err)
		assert.Equal(t, id.UnixMilli, uint64(1526985054069))
		assert.Equal(t, id.Seq, uint64(1))
	}
	a, err := red.NewStreamID("1526985054069-0")
	assert.NoError(t, err)
	b, err := red.NewStreamID("1526985054069-1")
	assert.NoError(t, err)
	assert.Equal(t, a.Bytes(), []byte("1526985054069-0"))
	assert.Equal(t, b.Bytes(), []byte("1526985054069-1"))
	assert.Equal(t, a.String(), "1526985054069-0")
	assert.Equal(t, b.String(), "1526985054069-1")
	assert.Equal(t, red.StreamID{}.Bytes(), []byte("0-0"))
	assert.Equal(t, red.StreamID{}.String(), "0-0")
	{
		assert.Equal(t, a.After(a), false)
		assert.Equal(t, a == a, true)
		assert.Equal(t, a.After(b), false)

		assert.Equal(t, b.After(b), false)
		assert.Equal(t, b == b, true)
		assert.Equal(t, b.After(a), true)
	}
	{
		assert.Equal(t, a.Before(a), false)
		assert.Equal(t, a == a, true)
		assert.Equal(t, a.Before(b), true)

		assert.Equal(t, b.Before(b), false)
		assert.Equal(t, b == b, true)
		assert.Equal(t, b.Before(a), false)
	}
}