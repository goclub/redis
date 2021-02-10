# goclub/redis

> goclub/redis 用于解决一些 go redis 库的接口设计的过于粗糙，导致编写 redis 相关 go 代码像是在写动态语言。


## SET

> SET key value

[设置永不过期的字符串](./examples/strings/set/set_never_expire_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSetNeverExpire (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}

	//  SET example_set_never_expire hello
	err = red.SET{
		Key: "example_set_never_expire",
		Value: "hello",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	log.Print(red.GET{Key: "example_set_never_expire"}.Do(ctx, radixClient))
	// "hello" true nil

	// 如果你没有传入 NeverExpire: true ，则会返回一个错误提醒你可能忘记配置过期时间
	err = red.SET{
		Key:"example_set_never_expire",
		Value: "hello",
	}.Do(ctx, radixClient)
	assert.Error(t, err, red.ErrSetForgetTimeToLive)
}

```

> EX PX 

[设置过期的字符串(基于 time.Duration)](./examples/strings/set/set_expire_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetExpire (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//  SET example_set abc PX 1000
	err = red.SET{
		Key: "example_set_expire",
		Value: "abc",
		// 只需配置 Expire 为 time.Duration，无需配置 EX PX ,goclub/redis 会将 Expire 自动转换为 PX
		Expire: time.Second * 1,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_expire"}.Do(ctx, radixClient))
	// "abc" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_expire"}.Do(ctx, radixClient))
	// "" false nil
}

```

> EXAT PXAT

[设置过期的字符串(基于 time.Time)](./examples/strings/set/set_expire_at_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetExpireAt (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//   SET example_set_expire_at nimoc PXAT timestamp-milliseconds
	err = red.SET{
		Key: "example_set_expire_at",
		Value: "nimoc",
		// 只需配置 ExpireAt 为 time.Time，无需配置 EXAT PXAT ,goclub/redis 会将 ExpireAt 自动转换为 PXAT
		ExpireAt: time.Now().Add(time.Second),
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_expire_at"}.Do(ctx, radixClient))
	// "nimoc" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_expire_at"}.Do(ctx, radixClient))
	// "" false nil
}

```

> KEEPTLL

[设置字符串且不修改过期时间](./examples/strings/set/set_keepttl_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"log"
	"testing"
	"time"
)

func TestSetKEEPTTL (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	//
	err = red.SET{
		Key: "example_set_keep_ttl",
		Value: "x",
		Expire: time.Second * 1,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "x" true nil
	radixClient.DebugOnce()
	// SET example_set_keep_ttl xyz KEEPTTL
	err = red.SET{
		Key: "example_set_keep_ttl",
		Value: "y",
		// 注意 KEEPTTL 是 6.0 才支持的功能
		KeepTTL:true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "y" true nil

	time.Sleep(time.Second)

	log.Print(red.GET{Key: "example_set_keep_ttl"}.Do(ctx, radixClient))
	// "" false nil
}

```

> NX

[设置不存在的key，如果存在返回 false](./examples/strings/set/set_nx_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetNX (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}
	_, err = red.DEL{Key: "example_set_nx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	//  SET example_set_nx hello1 NX
	ok, err := red.SETNX{
		Key: "example_set_nx",
		Value: "hello1",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第一次成功
	assert.Equal(t, ok, true)

	//  SET example_set_nx hello2 NX
	ok, err = red.SETNX{
		Key: "example_set_nx",
		Value: "hello2",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第二次因为 key 已存在所以失败
	assert.Equal(t, ok, false)

	// GET example_set_nx
	value, hasValue, err := red.GET{Key:"example_set_nx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	assert.Equal(t, value, "hello1")
	assert.Equal(t, hasValue, true)
}

```

> XX

[设置已存在的key,如果不存在返回 false](./examples/strings/set/set_xx_test.go)
```.go
package examples_strings_set

import (
	"context"
	red "github.com/goclub/redis"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetXX (t *testing.T) {
	ctx := context.Background()
	coreClient, err := (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil { panic(err) }
	defer coreClient.Close()
	radixClient := red.DriverRadixClient4{Core: coreClient}

	//  SET example_set_xx hello XX
	ok, err := red.SETXX{
		Key: "example_set_xx",
		Value: "hello1",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第一次失败，因为 key 不存在
	assert.Equal(t, ok, false)

	// 第二次 SET key value 设置值
	//  SET example_set_xx hello
	err = red.SET{
		Key: "example_set_xx",
		Value: "hello2",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}

	//  SET example_set_xx hello XX
	ok, err = red.SETXX{
		Key: "example_set_xx",
		Value: "hello3",
		NeverExpire: true,
	}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	// 第三次成功，因为 key 存在
	assert.Equal(t, ok, true)


	// GET example_set_xx
	value, hasValue, err := red.GET{Key:"example_set_xx"}.Do(ctx, radixClient) ; if err != nil {
		panic(err)
	}
	assert.Equal(t, value, "hello3")
	assert.Equal(t, hasValue, true)
}

```