package main

import (
	"bytes"
	"fmt"
)
const (
    CommandSET = "SET"
	CommandGET = "GET"
	CommandHELLO = "hello"
)
type Command interface {

}

type SetCommand struct {
	key, val []byte
}

type HelloCommand struct {
	value string
}
type GetCommand struct {
	key []byte
}

func respWriteMap(m map[string]string) []byte {
	buf := bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("+%s\r\n", k))
		buf.WriteString(fmt.Sprintf(":%s\r\n", v))
	}
	return buf.Bytes()
}