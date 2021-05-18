package red

import "context"

type EXISTS struct {
	Key string
	Keys []string
}
func (data EXISTS) Do(ctx context.Context, client Connecter) (existsCount uint64, err error) {
	args := []string{"EXISTS"}
	if data.Key != "" {
		data.Keys = []string{data.Key}
	}
	args = append(args, data.Keys...)
	var value int64
	value,_, err = client.DoIntegerReply(ctx, args) ; if err != nil {
		return
	}
	existsCount = uint64(value)
	return
}

