package red

import (
	"context"
	"time"
)
// SetLimiter 集合限制器
// 使用场景: 限制用户每天只能试读3个章节(如果不允许一天内反复试读相同章节则可以使用 IncrLimiter )
// 注意:
// 如果 key = free_trial:{userID} 			  Expire = 24h 是限制24小时
// 如果 key = free_trial:2022-01-01:{userID}   Expire = 24h 是限制每天
type SetLimiter struct {
	Key      string        `eg:"free_trial:2022-01-01:{userID}"`
	Member   string         `eg:"{chapterID}"`
	Expire   time.Duration `note:"有效期" eg:"time.Hour*24"`
	Maximum  uint64        `note:"最大限制" eg:"3"`
}

func (v SetLimiter) Do(ctx context.Context, client Connecter) (limited bool, err error) {
	// 涉及命令 SISMEMBER SCARD SADD PEXPIRE
	panic("TODO")
	return
}