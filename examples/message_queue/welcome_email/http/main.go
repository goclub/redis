package main

import (
	"context"
	xerr "github.com/goclub/error"
	xhttp "github.com/goclub/http"
	red "github.com/goclub/redis"
	examplesMQEmail "github.com/goclub/redis/examples/message_queue/welcome_email/email"
	exmapleMQ "github.com/goclub/redis/examples/message_queue/welcome_email/mesasge_queue"
	examplesMQUser "github.com/goclub/redis/examples/message_queue/welcome_email/user"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

func Handle(router *xhttp.Router, radixClient red.DriverRadixClient4) {
	router.HandleFunc(xhttp.Pattern{"GET", "/"}, func(c *xhttp.Context) (reject error) {
		query := c.Request.URL.Query()
		name := query.Get("name")
		if name == "" { name = "anonymous" + strconv.FormatInt(time.Now().Unix(), 10) }
		examplesMQEmail.SyncSendEmail(name)
		return c.WriteBytes([]byte("hello " + name +" (synchronous send email)"))
	})
	router.HandleFunc(xhttp.Pattern{"GET", "/use_message_queue"}, func(c *xhttp.Context) (reject error) {
		ctx := c.RequestContext()
		query := c.Request.URL.Query()
		name := query.Get("name")
		// 消息队列还有一个作用就是解耦，用户模块发布注册消息。
		// 邮件，短信等模块订阅。这样即使后续增加欢迎短信，用户模块也不需要修改代码。
		reject = examplesMQUser.PublishUserSignInMessage(ctx, radixClient, name) ; if reject != nil {
			return
		}
		if name == "" { name = "anonymous" + strconv.FormatInt(time.Now().Unix(), 10) }
		return c.WriteBytes([]byte("hello " + name +" (message queue send email)"))
	})
}
func main () {
	radixClient, err := exmapleMQ.ConnectRedis() ; if err != nil {
		panic(err)
	}

	router := xhttp.NewRouter(xhttp.RouterOption{
		OnCatchPanic: func(c *xhttp.Context, recoverValue interface{}) error {
			debug.PrintStack()
			log.Print(recoverValue)
			c.WriteStatusCode(500)
			return nil
		},
		OnCatchError: func(c *xhttp.Context, err error) error {
			var shouldRecord bool
			if reject, asReject := xerr.AsReject(err) ; asReject {
				shouldRecord = reject.ShouldRecord
				err := c.WriteBytes(reject.Response) ; if err != nil {
					return err
				}
			} else {
				shouldRecord = true
				c.WriteStatusCode(500)
				err := c.WriteBytes([]byte("server error,(unknown error)")) ; if err != nil {
					return err
				}
			}
			if shouldRecord {
				// 正式环境换成类似 sentry 的日志系统
				debug.PrintStack()
				log.Print(err)
			}
			return nil
		},
	})
	Handle(router, radixClient)
	addr := ":3000"
	serve := http.Server{
		Handler: router,
		Addr: addr,
	}
	log.Print("http://127.0.0.1" + addr)
	router.LogPatterns()
	go func() {
		listenErr := serve.ListenAndServe() ; if listenErr !=nil {
			if listenErr != http.ErrServerClosed {
				panic(listenErr)
			}
		}
	}()
	xhttp.GracefulClose(func() {
		log.Print("Shuting down server...")
		if err := serve.Shutdown(context.Background()); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
	})
}
