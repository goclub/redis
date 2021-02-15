package red

import (
	"context"
	"fmt"
	"strings"
)

func Command(ctx context.Context, client Client, valuePtr interface{}, args []string) (result Result, err error) {
	result, err = client.RedisCommand(ctx, valuePtr, args) ; if err != nil {
		if err.Error() == "ERR syntax error" {
			err = fmt.Errorf("%w: " + strings.Join(args, " "), err)
		}
	}
	return
}