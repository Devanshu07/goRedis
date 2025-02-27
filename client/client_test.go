package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewClient1(t *testing.T) {
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