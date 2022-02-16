package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTrigger_Do(t *testing.T) {
	for _, client := range Connecters {
		trigger_Do(t, client)
	}
}

func trigger_Do(t *testing.T, client Connecter) {
	ctx := context.TODO()
	do := func (namespace string) (bool, error) {
		return Trigger{
			Namespace: namespace,
			Interval: time.Second*1,
			Threshold: 3,
			Frequency: time.Second*1,
		}.Do(ctx, client)
	}
	_, err := DEL{
		Keys: []string{"trigger_alarm_1", "trigger_alarm_1:frequency"},
	}.Do(ctx, client)
	assert.NoError(t, err)
	// 第一次触发
	{
		triggered,err := do("trigger_alarm_1")
		assert.NoError(t, err)
		assert.Equal(t,triggered, false)
		value, isNil, err := GET{Key: "trigger_alarm_1"}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t,isNil, false)
		assert.Equal(t,value, "1")

	}
	// 第二次触发
	{
		triggered,err := do("trigger_alarm_1")
		assert.NoError(t, err)
		assert.Equal(t,triggered, false)

		value, isNil, err := GET{Key: "trigger_alarm_1"}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t,isNil, false)
		assert.Equal(t,value, "2")
	}
	// 第三次触发
	{
		triggered,err := do("trigger_alarm_1")
		assert.NoError(t, err)
		assert.Equal(t,triggered, true)

		value, isNil, err := GET{Key: "trigger_alarm_1"}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t,isNil, true)
		assert.Equal(t,value, "")
		result, err := PTTL{
			Key: "trigger_alarm_1:frequency",
		}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Less(t, result.TTL.Milliseconds(), int64(1000))
		assert.Greater(t, result.TTL.Milliseconds(), int64(500))
	}
	// 第456789次触发
	{
		// 4
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, false)
		}
		// 5
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, false)
		}
		// 6 不触发
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, false)
		}
		// sleep 1秒等待 frequency 过期
		time.Sleep(time.Second)
		// 7
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, false)
		}
		// 8
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, false)
		}
		// 9 触发
		{
			triggered,err := do("trigger_alarm_1")
			assert.NoError(t, err)
			assert.Equal(t,triggered, true)
		}
	}
}
