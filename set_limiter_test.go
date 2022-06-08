package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestSetLimiter_Do(t *testing.T) {
	for _, client := range Connecters {
		setLimiter_Do(t, client)
	}
}

func setLimiter_Do(t *testing.T, client Connecter) {
	ctx := context.TODO()
	_, err := DEL{
		Keys: []string{"set_Limiter_1"},
	}.Do(ctx, client)
	_, err = DEL{
		Keys: []string{"set_Limiter_2"},
	}.Do(ctx, client)
	assert.NoError(t, err)
	/* 第1组 第1个 第1次 */
	limited, err := SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "a",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	// 验证过期时间
	result, err := PTTL{
		Key: "set_Limiter_1",
	}.Do(ctx, client);assert.NoError(t, err)
	log.Print(result.TTL)
	tenSecond := time.Second*10
	assert.Less(t, result.TTL.Milliseconds(), tenSecond.Milliseconds())
	/* 第1组 第1个 第2次 */
	time.Sleep(time.Second)
	limited, err = SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "a",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第1组 第2个 */
	time.Sleep(time.Second)
	limited, err = SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "b",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第1组 第3个 */
	time.Sleep(time.Second)
	limited, err = SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "c",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第1组 第4个 */
	time.Sleep(time.Second)
	limited, err = SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "d",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, true)
	/* 第2组 */
	time.Sleep(time.Second)
	limited, err = SetLimiter{
		Key:     "set_Limiter_2",
		Member:  "a",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
	/* 第1组 过期后 */
	time.Sleep(time.Second*7) // 等待过期
	limited, err = SetLimiter{
		Key:     "set_Limiter_1",
		Member:  "e",
		Expire:  time.Second * 10,
		Maximum: 3,
	}.Do(ctx, client);assert.NoError(t, err)
	assert.Equal(t, limited, false)
}
