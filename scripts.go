package red

import "context"

type Script struct {
	Keys []string
	Argv []string
	Script string
}
func (data Script) Do(ctx context.Context, client Connecter) (reply interface{}, isNil bool, err error) {
	return client.Eval(ctx, data)
}