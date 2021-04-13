# 时间频率限制

## 需求

在广告业务，限制某个广告主对于一个用户在限定时间内内只曝光一次。

>  限制的时间不确定，是可以根据后台配置控制。可以是一小时或一天或两天

## sql 实现

如果使用 sql 实现，则需要记录每一次用户曝光广告记录下来。通过 sql 查询

```sql
SELECT 
    time 
FORM
    advertising_record
WHERE 
    user_id = "u1"
ORDER BY
    time DESC 
LIMIT 1
``` 

获取最近一次曝光的时间，将曝光时间与当前时间计算出时间间隔。
判断时间间隔是否大于10秒，如果大于则曝光广告，反之不曝光广告。


伪代码如下:

```
function limit(userID string) {
    record = queryLatestRecotd()
    if record == nil {
        insertRecord(userID, now())
        return "pass"
    }
    limitDuration = 10s
    if now() - record > limitDuration {
        insertRecord(userID, now())
        return "pass"
    }
    return "limited"
}
```

虽然这样能实现，但是性能不好，每次都要查询sql,并且查询->判断->插入三个操作不满足原子性，在高并发的场景下会导致多次曝光。

> 注意在曝光广告的场景下由高并发导致的曝光是理论上会出现，实际上曝光广告的触发不会高并发,所以是不需要解决不满足原子性带来的数据不一致。
> 但从学习的角度，我们要让代码符合原子性


# 简单但会误判的方法

`SET key value PX $limitDurationMilli NX`

如果响应 OK 则限制，否则不限制。

但是如果后台修改了配置。把时间间隔从2天改为了1天，这个配置对于曝光过广告的用户必须一天后才能生效。

如果业务上允许出现部分误判情况，那就没问题。   


# 无误判的方法

想要做到无误判，需要在 redis 记录曝光时间。

但是没有必要在 redis 存储用户每一次曝光记录。只需要存储一条记录，并且不停的更新这条记录。

redis 实现代码: [basic/main.go](./basic/main.go)

# 扩展

如果后续代码执行或者网络io出现错误需要考虑是否应当调用 ClearLimit 
要注意 Limited() -> someAction() -> ClearLimit() 是不满足原子性的,
这需要在具体的业务场景中考虑不满足原子性造成的后果是否能接受，这里我就不深入讨论了
```
if Limited() {
    return
}
err = someAction() ; if err != nil {
    err = ClearLimit(ctx, client, userID) ; if err != nil {
        panic(err)
    }
}

func ClearLimit (ctx context.Context, client red.Connecter, userID string) (err error) {
	recordTimeKey := "example_tfl:" + userID
	_, err = red.DEL{Key: recordTimeKey}.Do(ctx, client)
	return
}
```