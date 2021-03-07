# goclub/redis

> goclub/redis 用于解决一些 go redis 库的接口设计的过于粗糙，导致编写 redis 相关 go 代码像是在写动态语言。

## strings

### SET

> SET key value

[设置永不过期的字符串|embed](./examples/strings/set/set_never_expire_test.go)

> EX PX 

[设置过期的字符串(基于 time.Duration)|embed](./examples/strings/set/set_expire_test.go)

> EXAT PXAT

[设置过期的字符串(基于 time.Time)|embed](./examples/strings/set/set_expire_at_test.go)

> KEEPTLL

[设置字符串且不修改过期时间|embed](./examples/strings/set/set_keepttl_test.go)

> NX

[设置不存在的key，如果存在返回 false|embed](./examples/strings/set/set_nx_test.go)

> XX

[设置已存在的key,如果不存在返回 false|embed](./examples/strings/set/set_xx_test.go)

### GET

> GET key

[GET key|embed](./examples/strings/get/get_test.go)

### APPEND

[APPEND key value |embed](./examples/strings/append/append_test.go)



