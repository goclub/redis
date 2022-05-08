package example_test

import (
	"context"
	"log"
	"testing"
)

func TestExampleDoIntegerReply (t *testing.T) {
	ExampleDoIntegerReply()
}
func ExampleDoIntegerReply()  {
    ctx := context.Background()
	err := func() (err error){
		client, err := NewClient(ctx) ; if err != nil {
		    return
		}
		replyInt64, err := client.DoIntegerReplyWithoutNil(
			ctx,
			[]string{
				"PFADD", "example_hhl", "a",
			},
		) ; if err != nil {
			return
		}
		log.Print("PFADD reply is ", replyInt64)
		// 因为PFADD 的返回是 integer-reply 所以使用 DoIntegerReplyWithoutNil
		// 如果返回可能有 nil 则使用 DoIntegerReply
		// 还有 DoStringReply DoArrayIntegerReply DoArrayStringReply 等方法可以使用
		// 它们对应了 redis RESP 的各种返回值
		return
	}() ; if err != nil {
	    log.Printf("%+v",err)
	}
}