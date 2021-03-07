package examples_append_get

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	radix4 "github.com/redis-driver/mediocregopher-radix-v4"
	"log"
	"testing"
)

func TestAppend(t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	client := radix4.Client{Core: coreClient}
	length, err :=  red.APPEND{Key: "example_get"}.Do(ctx, client) ; if err != nil {
		panic(err)
	}
	log.Print(length)
}