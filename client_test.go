package red_test

import (
	"context"
	"github.com/mediocregopher/radix/v4"
)

var radixClient red.Radix4Client
func init () {
	ctx := context.Background()
	client, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
		panic(err)
	}
	radixClient = red.Radix4Client{client}
}
// func TestNewClient(t *testing.T) {
// 	ctx := context.Background()
// 	client := Test{}
// 	_, err := red.SET{
// 		Key: "name",
// 		Value: "tim",
// 		Expire: time.Minute,
// 		NX: true,
// 	}.Do(ctx, client) ; if err != nil {
// 		panic(err)
// 	}
// 	value, hasValue, err := red.GET{
// 		Key: "name",
// 	}.Do(ctx, client) ; if err != nil {
// 		panic(err)
// 	}
// 	log.Print(value, hasValue)
// }