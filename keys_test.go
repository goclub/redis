package red

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDel(t *testing.T) {
	for _, client := range Connecters {
		redisDel(t, client)
	}
}
func redisDel(t *testing.T, client Connecter) {
	ctx := context.TODO()
	key := "del"
	key2 := "del2"
	{// SET key v
		_, _, err := client.DoStringReply(ctx, []string{"SET", key, "v"}) ; assert.NoError(t, err)
	}
	{// DEL key
		delCount, err := DEL{Key: key}.Do(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t, delCount, uint64(1))
		reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key})
		assert.NoError(t, err)
		assert.Equal(t, reply, "")
		assert.Equal(t, isNil, true)}
	{
		// MSET key "nimo" key2 "nico"
		_, _, err := client.DoStringReply(ctx, []string{"MSET", key, "nimo", key2, "nico"}) ; assert.NoError(t, err)
	}
	{// DEL key key2
		delCount, err := DEL{Keys: []string{key, key2}}.Do(ctx, client)
		assert.Equal(t, delCount, uint64(2))
		assert.NoError(t, err)
		// GET key
		{
			reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key})
			assert.Equal(t, reply, "")
			assert.Equal(t, isNil, true)
			assert.NoError(t, err)
		}
		// GET key2
		{
			reply, isNil, err := client.DoStringReply(ctx, []string{"GET", key2})
			assert.Equal(t, reply, "")
			assert.Equal(t, isNil, true)
			assert.NoError(t, err)
		}
	}
}
