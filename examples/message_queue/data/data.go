package exmaplesMQData

type UserSignIn struct {
	UserID string
}
func (UserSignIn) StreamKey() string {
	return "mq_user_sign_in"
}