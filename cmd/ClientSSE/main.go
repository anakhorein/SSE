package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/r3labs/sse/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	batchMin, batchMax := 2, 10
	batchSize := rand.Intn(batchMax-batchMin) + batchMin
	client := sse.NewClient(fmt.Sprintf("http://localhost:8080/task?batchSize=%d", batchSize))

	client.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("goroutine paniqued: ", r)
				}
			}()
			//fmt.Println(msg.Data)
			buf := bytes.NewBuffer(msg.Data)
			dec := gob.NewDecoder(buf)
			var message msg1
			err := dec.Decode(&message)
			if err != nil {
				log.Fatal("decode error:", err)
			}

			var i = message.Period
			println(i)
			if i > 900 {
				os.Exit(0)
			}
			if i > 800 {
				panic(i)
			}
			time.Sleep(time.Duration(i) * time.Millisecond)

		}()

	})
}

type msg1 struct {
	Id     string // https://github.com/rs/xid
	Period uint64 // 1-1000 random
}
