---
permalink: /
sidebarBasedOnContent: true
---

# goclub/redis

[![Go Reference](https://pkg.go.dev/badge/github.com/goclub/redis.svg)](https://pkg.go.dev/github.com/goclub/redis)

## 易用性

go 社区有很多 [redis 库](https://redis.io/docs/clients/#go),
但是大多数库在 redis 返回 nil 时候用 error 表示.例如使用GET方法:

```go
cmd := client.Get(ctx, "name")
err = cmd.Err()
isNil := false
if err != nil {
    if errors.Is(err, redis.Nil) {
        isNil = true
    } else {
        return err
    }
}
if isNil {
    // do some
} else {
    log.Print(cmd.String())
}
```

代码写的复杂,心很累.


goclub/redis 的接口风格是

```go
value, isNil, err := red.GET{Key:key}.Do(ctx, client) ; if err != nil {
    return
}
if isNil {
	// do some
} else {
	log.Print(value)
}
```


## 自由

redis 的核心是 [RESP](https://redis.io/docs/reference/protocol-spec/) 协议. goclub/redis 提供以下接口对应 RESP.
```go
// Connecter RESP
type Connecter interface {
	DoStringReply(ctx context.Context, args []string) (reply string, isNil bool, err error)
	DoStringReplyWithoutNil(ctx context.Context, args []string) (reply string, err error)
	DoIntegerReply(ctx context.Context, args []string) (reply int64, isNil bool, err error)
	DoIntegerReplyWithoutNil(ctx context.Context, args []string) (reply int64, err error)
	DoArrayIntegerReply(ctx context.Context, args []string) (reply []OptionInt64, err error)
	DoArrayStringReply(ctx context.Context, args []string) (reply []OptionString, err error)
	Eval(ctx context.Context, script Script) (reply Reply, isNil bool, err error)
	EvalWithoutNil(ctx context.Context, script Script) (reply Reply, err error)
}
```

你可以自由的使用任何命令.

```go
replyInt64, err = client.DoIntegerReplyWithoutNil(
	ctx, 
	[]string{
		"HSET", key, "name", "goclub", "age", "18",
	})
if err != nil {
    return
}
```

你也可以查看 `red.API` 查看 goclub/redis 对哪些方法进行了封装

---

像是 `RESP Simple Strings` 这类操作你可以根据 redis 是否会返回 nil 去调用 `DoStringReply` 或 `DoStringReplyWithoutNil`

脚本则使用  `Eval` 或者`EvalWithoutNil`

```go
script := `
if redis.call("GET", KEYS[1]) == ARGV[1]
then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`
reply, err := client.EvalWithoutNil(ctx, Script{
    KEYS: []string{key},
    ARGV: []string{value},
    Script: script,
}) ; if err != nil {
    return
}
/*
你可以根据 script 内容调用一下 reply 方法
reply.Int64()
reply.Uint64()
reply.String()
reply.Int64Slice()
reply.StringSlice()
*/
```

## 包容

你可以将 goclub/redis 理解为是一个驱动库,是一个"壳".
goclub/redis 能通过接口适配 go 社区的说有 redis 库.

> [goredis](https://redis.uptrace.dev/) 是非常流行的库, goclub/redis 默认支持了 [goredisv8](./goredisv8.go)


## 直接来吧

![](./start.jpg)

[连接redis    | NewClient](./example/internal/new_client_test.go?embed)

[字符串增删改查 | StringsCRUD](./example/internal/strings_crud_test.go?embed)

[eval执行脚本 | Eval](./example/internal/eval_test.go?embed)

[直接写命令 |l DoIntegerReply ](./example/internal/do_interger_reply_test.go?embed)


## 实用方法

goclub/redis 还基于redis实现了了一些实用的方法,帮助开发人员提高开发效率.防止重复造轮子.

### Trigger

> 触发器

```go
// Trigger 触发器
// 每5分钟出现3次则触发,但是10分钟内只触发一次
func exmaple () {
	triggered, err := red.Trigger{
		Namespace: "pay_fail_alarm",
		Interval: time.Minute*5,
		Threshold: 3,
		Frequency: time.Minute*5,
	}.Do(ctx, client) ; if err != nil {
	    return
	}
	if triggered {
		// do some
	}
}
```

### Mutex

> 互斥锁

```go
mutex := red.Mutex{
    Key: key,
	// Expire 控制锁定时间
    Expire: time.Millisecond*100,
	// Retry 当锁被占用时进入循环重试(此时会堵塞)
	// Retry.Times 重试上锁多少次
	// Retry.Interval 重试上锁间隔
    Retry: red.Retry{
        Times: 3,
		Interval:time.Millisecond*100,
    },
}
lockSuccess, unlock, err := mutex.Lock(context.TODO(), client) ; if err != nil {
    // 锁失败
    return
}
if lockSuccess == false {
    // 锁被占用
    return
}
// 处理某些业务
err = unlock(context.TODO()) ; if err != nil {
    *mutexCount--
	// 解锁失败
	log.Printf("%+v", err)
    return
}
// 解锁成功
```

### IncrLimiter

> 递增限制器

```go
alarm_1 := IncrLimiter{
    Namespace: "incr_limiter_alarm_1",
    Expire:    time.Second * 10,
    Maximum:   3,
}
/* 第1次 */
limited, err := alarm_1.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第2次 */
limited, err := alarm_1.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第3次 */
limited, err := alarm_1.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第4次 */
limited, err := alarm_1.Do(ctx, client) ; if err != nil {
    return
} // limited = true
```


### SetLimiter

> 设值限制器

```go
/* 第1次 */
limited, err := SetLimiter{
    Namespace: "set_limiter_alarm_1",
    Member:    "a"
    Expire:    time.Second * 10,
    Maximum:   3,
}.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第1次 重复 */
limited, err := SetLimiter{
    Namespace: "set_limiter_alarm_1",
    Member:    "a"
    Expire:    time.Second * 10,
    Maximum:   3,
}.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第2次 */
limited, err := SetLimiter{
    Namespace: "set_limiter_alarm_1",
    Member:    "b"
    Expire:    time.Second * 10,
    Maximum:   3,
}.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第3次 */
limited, err := SetLimiter{
    Namespace: "set_limiter_alarm_1",
    Member:    "c"
    Expire:    time.Second * 10,
    Maximum:   3,
}.Do(ctx, client) ; if err != nil {
    return
} // limited = false

/* 第4次 */
limited, err := SetLimiter{
    Namespace: "set_limiter_alarm_1",
    Member:    "d"
    Expire:    time.Second * 10,
    Maximum:   3,
}.Do(ctx, client) ; if err != nil {
    return
} // limited = true
```