package examples_bit_get_test

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	radix4 "github.com/redis-driver/mediocregopher-radix-v4"
	"testing"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	client := radix4.Client{Core: coreClient}
}
