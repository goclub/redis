package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHdel(t *testing.T) {
	for _, client := range Connecters {
		redisHdel(t, client)
	}
}

func redisHdel(t *testing.T, client Connecter) {
	func() struct{} {
		// -------------
		ctx := context.TODO()
		key := "hdel"
		_, err := DEL{Key: key}.Do(ctx, client) ; assert.NoError(t, err)
		{
			_, err = client.DoIntegerReplyWithoutNil(ctx, []string{"HSET", key, "name", "goclub"})
			assert.NoError(t, err)
			delTotal, err := HDEL{
				Key:   key,
				Field: []string{"name"},
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Equal(t, delTotal, uint64(1))
		}
		{
			delTotal, err := HDEL{
				Key:   key,
				Field: []string{"name"},
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Equal(t, delTotal, uint64(0))
		}
		{
			_, err = client.DoIntegerReplyWithoutNil(ctx, []string{"HSET", key, "name", "goclub", "age", "18"})
			assert.NoError(t, err)
			delTotal, err := HDEL{
				Key:   key,
				Field: []string{"name", "age"},
			}.Do(ctx, client) ; assert.NoError(t, err)
			assert.Equal(t, delTotal, uint64(2))
		}
		// -------------
		return struct{}{}
	}()
}
