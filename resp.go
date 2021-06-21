package red

func ParseArrayIntegerReply(reply interface{}) (arrayIntegerReply []OptionInt64) {
	values := reply.([]interface{})
	for _, v := range values {
		if v == nil {
			arrayIntegerReply = append(arrayIntegerReply, OptionInt64{})
		} else {
			arrayIntegerReply = append(arrayIntegerReply, NewOptionInt64(v.(int64)))
		}
	}
	return
}
func ParseArrayStringReply(reply interface{}) (arrayStringReply []OptionString) {
	values := reply.([]interface{})
	for _, v := range values {
		if v == nil {
			arrayStringReply = append(arrayStringReply, OptionString{})
		} else {
			arrayStringReply = append(arrayStringReply, NewOptionString(v.(string)))
		}
	}
	return
}