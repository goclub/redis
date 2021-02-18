package red_test

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	radix4 "github.com/redis-driver/mediocregopher-radix-v4"
)

var radixClient radix4.Client
func init () {
	ctx := context.Background()
	client, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
		panic(err)
	}
	// 测试用场景，所以省略 close
	// client.Close()
	radixClient = radix4.Client{Core: client}
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