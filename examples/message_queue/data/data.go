package exampleMQData

type UserSignIn struct {
	UserID string `red:"user_id"`
}
func (UserSignIn) StreamKey() string {
	return "mq_user_sign_in"
}