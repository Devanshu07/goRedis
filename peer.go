package main

import (
	"fmt"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn net.Conn // Underlying network connection (e.g., TCP socket)
	msgCh chan Message
}

func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn: conn,
		msgCh: msgCh,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

        if v.Type() == resp.Array {
            for _, value := range v.Array() {
                switch value.String() {
				case CommandGET:
					if(len(v.Array()) != 2) {
						return fmt.Errorf("invalid number of variables for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					fmt.Println("Got GET cmd %+v\n", cmd)
				case CommandSET:
					if(len(v.Array()) != 3) {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					fmt.Println("Got SET cmd %+v\n", cmd)
				}
            }
        }
	}
	return nil
}