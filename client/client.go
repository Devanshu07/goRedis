package client

import (
	"bytes"
	"context"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		addr : addr,
		conn: conn,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error{
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"), 
		resp.StringValue(key), 
		resp.StringValue(val),
	})
	_, err = conn.Write(buf.Bytes())
	return err
}