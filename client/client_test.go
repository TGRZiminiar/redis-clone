package client

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestOfficialRedisClient(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// fmt.Println("this is rdb => ", rdb)

	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		t.Fatal(err)
	}

	val, err := rdb.Get(ctx, "foo").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}

func TestNewClient1(t *testing.T) {
	client, err := NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	fmt.Println("SET => ", "bar")
	if err := client.Set(context.Background(), "foo", "hello"); err != nil {
		log.Fatal(err)
	}
	val, err := client.Get(context.Background(), "foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GET => ", val)

}

func TestNewClient(t *testing.T) {
	client, err := NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("SET => ", fmt.Sprintf("bar %d", i))
		if err := client.Set(context.Background(), fmt.Sprintf("foo %d", i), fmt.Sprintf("bar %d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get(context.Background(), fmt.Sprintf("foo %d", i))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("GET => ", val)

	}

}
