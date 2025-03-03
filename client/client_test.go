package client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)


func TestNewClientRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:5001",
		Password: "",
		DB: 0,
	})
	fmt.Println(rdb)
	fmt.Println("this is working fine")

	// err := rdb.Set(context.Background(), "key", "value", 0).ERR()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(context.TODO(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(val)
}
func TestNewClient1(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	time.Sleep(time.Second)
	if err := c.Set(context.TODO(), "foo", 69); err != nil {
		log.Fatal(err)
	}
	val, err := c.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	n, _ := strconv.Atoi(val)
	fmt.Println(n)
	fmt.Println("GET => ", val)
}

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Set(context.TODO(), "foo", "1"); err != nil {
		log.Fatal(err)
	}
	val, err := c.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}