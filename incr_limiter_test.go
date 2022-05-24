package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIncrLimiter_Do(t *testing.T) {
	ctx := context.TODO()
	for i, client := range Connecters {
		if i == 0 {
			_, err := DEL{
				Keys: []string{"incr_limiter_alarm_1"},
			}.Do(ctx, client)
			assert.NoError(t, err)
		}
		incrLimiter_Do(t, client)
	}
}

func incrLimiter_Do(t *testing.T, client Connecter) {
	ctx := context.TODO()
	alarm_1 := IncrLimiter{
		Key:     "incr_limiter_alarm_1",
		Expire:  time.Second * 10,
		Maximum: 3,
	}
	/* 第1次 */
	limited, err := alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	// 验证过期时间
	result, err := PTTL{
		Key: "incr_limiter_alarm_1",
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Less(t, result.TTL.Milliseconds(), alarm_1.Expire.Milliseconds())
	/* 第2次 */
	time.Sleep(time.Second)
	limited, err = alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第3次 */
	time.Sleep(time.Second)
	limited, err = alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第4次 */
	time.Sleep(time.Second)
	limited, err = alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, true)
	/* 第5次 */
	time.Sleep(time.Second)
	limited, err = alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, true)
	/* 第6次 */
	time.Sleep(time.Second*7) // 等待过期
	limited, err = alarm_1.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	// 验证ppl 小于20分钟
}
