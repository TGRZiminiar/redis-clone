package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"tgrziminiar/redisclone/client"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestServer(t *testing.T) {

	server := NewServer(&Config{})
	go func() {

		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// fmt.Println("this is rdb => ", rdb)
	// _ = ctx
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}

func TestRESP(t *testing.T) {
	in := map[string]string{
		"server":  "redis",
		"version": "6.0",
	}

	out := respWriteMap(in)
	fmt.Println(string(out))

}

func TestServerWithMultiClient(t *testing.T) {

	server := NewServer(&Config{})
	go func() {

		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	nClinets := 10
	wg := sync.WaitGroup{}
	wg.Add(nClinets)
	for i := 0; i < nClinets; i++ {
		go func(it int) {
			client, err := client.NewClient("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer client.Close()

			key := fmt.Sprintf("client_foo_%d", it)
			value := fmt.Sprintf("client_bar_%d", it)
			if err := client.Set(context.Background(), key, value); err != nil {
				log.Fatal(err)
			}
			val, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("client %d got this val back => %s \n", it, val)
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("expeceted 0 peers but got %d", len(server.peers))
	}

}
