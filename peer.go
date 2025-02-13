package main

import (
	"net"
)

type Peer struct {
	conn net.Conn // Underlying network connection (e.g., TCP socket)
	msgCh chan []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn: conn,
		msgCh: msgCh,
	}
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
		p.msgCh <- msgBuf

	}
}