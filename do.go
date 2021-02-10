package red

import (
	"context"
	"fmt"
	"strings"
)

func Command(ctx context.Context, doer Doer, valuePtr interface{}, args []string) (result Result, err error) {
	result, err = doer.RedisCommand(ctx, valuePtr, args) ; if err != nil {
		if err.Error() == "ERR syntax error" {
			err = fmt.Errorf("%w: " + strings.Join(args, " "), err)
		}
	}
	return
}