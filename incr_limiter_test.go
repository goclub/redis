package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIncrLimiter_Do(t *testing.T) {
	for _, client := range Connecters {
		incrLimiter_Do(t, client)
	}
}

func incrLimiter_Do(t *testing.T, client Connecter) {
	ctx := context.TODO()
	{
		key := "incr_limiter_alarm_1"
		_, err := DEL{
			Keys: []string{key},
		}.Do(ctx, client)
		assert.NoError(t, err)
		v := IncrLimiter{
			Key:     key,
			Expire:  time.Second * 10,
			Maximum: 3,
		}
		/* 第1次 */
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, false)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second)
			assert.Less(t, result.TTL, v.Expire)
		}
		/* 第2次 */
		time.Sleep(time.Second)
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, false)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second*2)
			assert.Less(t, result.TTL, v.Expire-time.Second)
		}
		/* 第3次 */
		time.Sleep(time.Second)
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, false)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second*3)
			assert.Less(t, result.TTL, v.Expire-time.Second*2)
		}
		/* 第4次 */
		time.Sleep(time.Second)
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, true)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second*4)
			assert.Less(t, result.TTL, v.Expire-time.Second*3)
		}
	}
	{
		key := "incr_limiter_alarm_2"
		_, err := DEL{
			Keys: []string{key},
		}.Do(ctx, client) ; assert.NoError(t, err)
		v := IncrLimiter{
			Key:     key,
			Expire:  time.Second * 10,
			Maximum: 3,
			Increment: 2,
		}
		/* 第1次 */
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, false)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client);assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second)
			assert.Less(t, result.TTL, v.Expire)
		}
		// 验证值
		{
			value, isNil, err := GET{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Equal(t,isNil, false)
			assert.Equal(t,value, "2")
		}
		/* 第2次 */
		{
			limited, err := v.Do(ctx, client);assert.NoError(t, err)
			assert.Equal(t, limited, true)
		}
		// 验证过期时间
		{
			result, err := PTTL{
				Key: key,
			}.Do(ctx, client);assert.NoError(t, err)
			assert.Greater(t, result.TTL, v.Expire-time.Second)
			assert.Less(t, result.TTL, v.Expire)
		}
		// 验证值
		{
			value, isNil, err := GET{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Equal(t,isNil, false)
			assert.Equal(t,value, "2")
		}
	}
}
