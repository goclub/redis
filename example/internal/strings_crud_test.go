package example_test

import (
	"context"
	red "github.com/goclub/redis"
	"log"
	"testing"
	"time"
)

func TestStringsCRUD(t *testing.T) {
	err := func() (err error){
		ctx := context.Background()
		client, err := NewClient(ctx) ; if err != nil {
			return
		}
		key := "example_strings_crud"

		// ----------------
		log.Print("DEL")
		delTotal, err := red.DEL{
			Key: key,
		}.Do(ctx, client) ; if err != nil {
		    return
		}
		log.Printf("删除%d个key", delTotal)
		// ----------------
		log.Print("GET")
		getReply, getIsNil, err := red.GET{
			Key: key,
		}.Do(ctx, client) ; if err != nil {
			return
		}
		if getIsNil {
			log.Print(key + " 为空")
		} else {
			log.Print(key + "的内容是", getReply)
		}

		// ----------------
		log.Print("SET EX")
		setExReply, setExIsNil,  err := red.SET{
			Key: key,
			Value: "a",
			// 20秒有效期
			Expire: time.Second * 20,
		}.Do(ctx, client) ; if err != nil {
		    return
		}
		// set ex 不会 返回 nil
		_ = setExIsNil
		// set ex 返回的字符串固定是 OK 所以忽略
		_ = setExReply

		// ----------------
		log.Print("SET NX")
		setNxReply, setNxIsNil, err := red.SET{
			Key: key,
			Value: "b",
			Expire: time.Second * 20,
			NX:  true,
		}.Do(ctx, client) ; if err != nil {
		    return
		}
		if setNxIsNil {
			log.Print(key + " 已经存在值,因为配置了 NX 所以值未改变")
		} else {
			log.Print(key + "set nx 成功")
		}
		// set nx 返回的字符串固定是 OK 所以忽略
		_ = setNxReply

		// ----------------
		log.Print("GET")
		getReply, getIsNil, err = red.GET{Key: key}.Do(ctx, client) ; if err != nil {
		    return
		}
		if getIsNil {
			log.Print(key + " 为空")
		} else {
			log.Print(key + "的内容是", getReply)
		}
		return
	}() ; if err != nil {
		log.Printf("%+v",err)
	}
}