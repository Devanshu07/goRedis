package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
)
const (
    CommandSET = "SET"
)
type Command interface {

}

type SetCommand struct {
	key, val []byte
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

        if v.Type() == resp.Array {
            for _, value := range v.Array() {
                switch value.String() {
				case CommandSET:
					fmt.Println(len(v.Array()))
					if(len(v.Array()) != 3) {
						return nil, fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					return cmd, nil
				}
            }
        }
		return nil, fmt.Errorf("invalid or unknown command received: %s", raw)
	}
    return nil, fmt.Errorf("invalid or unknown command received: %s", raw)
}