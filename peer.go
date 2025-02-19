package main

import (
	"net"
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
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err 
		}
		// fmt.Println(string(buf[:n]))
		// fmt.Println(len(buf[:n]))
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgCh <- Message{
			data: msgBuf,
			peer: p,
		}

	}
}