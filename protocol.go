package main

import (
	"fmt"
)
const (
    CommandSET = "SET"
	CommandGET = "GET"
)
type Command interface {

}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

func parseCommand(raw string) (Command, error) {
	// rd := resp.NewReader(bytes.NewBufferString(raw))
	
    return nil, fmt.Errorf("invalid or unknown command received: %s", raw)
}