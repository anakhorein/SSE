package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/r3labs/sse/v2"
	"github.com/rs/xid"
	bolt "go.etcd.io/bbolt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {

	// Open the 'my.db' data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("messages.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Messages"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	server := sse.New()
	server.CreateStream("messages")

	// Create a new Mux and set the handler
	mux := http.NewServeMux()

	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
			return
		}()

		batchSize, err := strconv.Atoi(r.URL.Query().Get("batchSize"))
		if err != nil {
			//panic(err)
			batchSize = 1
		}

		go getMessages(server, db, batchSize)
		server.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", mux)
}

func getMessages(server *sse.Server, db *bolt.DB, batchSize int) {
	println(batchSize)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < batchSize; i++ {

		guid := xid.New()

		var msgBytes bytes.Buffer
		enc := gob.NewEncoder(&msgBytes)
		err := enc.Encode(msg{Id: guid.String(), Period: uint64(r.Int63n(1000))})
		if err != nil {
			log.Fatal("encode error:", err)
		}
		if err != nil {
			print(err)
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Messages"))
			err := b.Put(guid.Bytes(), msgBytes.Bytes())
			return err
		})
		server.Publish("messages", &sse.Event{
			Data: msgBytes.Bytes(),
		})

	}
}

type msg struct {
	Id string // https://github.com/rs/xid

	Period uint64 // 1-1000 random

}
