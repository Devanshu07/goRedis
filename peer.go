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
	delCh chan *Peer
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn: conn,
		msgCh: msgCh,
		delCh: delCh,
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
			p.delCh <- p
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

					p.msgCh <- Message{
						cmd: cmd,
						peer: p,
					}

				case CommandSET:
					if(len(v.Array()) != 3) {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					
					p.msgCh <- Message{
						cmd: cmd,
						peer: p,
					}
				default:
					fmt.Println(v.Array())
				}
            }
        }
	}
	return nil
}