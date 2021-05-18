package red_test

import (
	"context"
	red "github.com/goclub/redis"
	"log"
)

func ExampleAPPEND_Do() {
	ctx := context.Background()
	// DEL mykey
	_, err := red.DEL{Key: "mykey"}.Do(ctx, exampleClient)
	if err != nil { log.Print(err);return }

	// EXISTS mykey
	existsCount, err := red.EXISTS{Key: "mykey"}.Do(ctx, exampleClient)
	if err != nil { log.Print(err);return }
	log.Print("existsCount", existsCount)

	// APPEND mykey "Hello"
	newLength, err :=  red.APPEND{Key: "mykey", Value: "Hello"}.Do(ctx, exampleClient)
	if err != nil { log.Print(err);return }
	log.Print("newLength", newLength)

	// APPEND mykey " World"
	newLength, err =  red.APPEND{Key: "mykey", Value: "World"}.Do(ctx, exampleClient)
	if err != nil { log.Print(err);return }
	log.Print("newLength", newLength)

	// GET mykey
	log.Print(
		red.GET{Key: "mykey"}.Do(ctx, exampleClient),
	)

}
