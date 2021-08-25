# goclub/redis

## 启动

```go
package main

import (
	red "github.com/goclub/redis"
	redis "github.com/go-redis/redis/v8"
	"context"
)

func main() {
	ctx := context.Background()
    exampleClient := red.GoRedisV8{
        Core: redis.NewClient(&redis.Options{}),
    }
    err := exampleClient.Core.Ping(ctx).Err() ; if err != nil {
        panic(err)
    }
    // red.GET{Key: "mykey"}.Do(ctx, exampleClient)
    // red.SET{Key: "mykey", Value: "nimo"}.Do(ctx, exampleClient)
    // red.DEL{Key: "mykey"}.Do(ctx, exampleClient)
}
```
