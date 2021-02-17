package exmapleMQ

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
)

type UserSignInMessage struct {
	Name string `red:"name"`
}
func (UserSignInMessage) StreamKey() string {
	return "mq_user_sign_in"
}

func ConnectRedis() (red.DriverRadixClient4, error) {
	ctx := context.Background()
	client, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
		return red.DriverRadixClient4{}, err
	}
	return red.DriverRadixClient4{Core: client}, err
}
