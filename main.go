package main

import (
	"context"
	"fmt"
	"goredis/client"
	"log"
	"log/slog"
	"net"
	"time"
)

const defaultListenAddr = ":5001" // Default port for the server
type Config struct {
	ListenAddr string // Configuration for the server address

}
type Server struct {
	Config 				 // Embeds the Config struct (inherits its fields)
	peers map[*Peer]bool // Tracks active peers (clients)
	ln net.Listener		 // TCP listener
	addPeerCh chan *Peer // Channel to add new peers
	quitCh chan struct{} // Channel to signal shutdown
	msgCh chan []byte
	kv *KV
}



func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr=defaultListenAddr // Use default port if none specified
	}
	return &Server{
		Config: cfg,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
		msgCh: make(chan []byte),
		kv: NewKV(),
	}
}

func (s *Server) Start() error{  // Start TCP listener
	ln, err :=net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return err
	}
	s.ln=ln
	go s.loop() // Start the event loop in a goroutine

	slog.Info("Server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop() // Block and accept connections
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <- s.msgCh:
			if err := s.handleRawMessage(rawMsg); err !=nil {
				slog.Error("Raw Message error", "err", err)
			} 
		case <- s.quitCh:  // Shutdown signal
			return
		case peer := <- s.addPeerCh: // Add a new peer
			s.peers[peer] = true
		}
	}
}

func(s *Server) acceptLoop() error{
	for {
		conn, err := s.ln.Accept()  // Block until a connection arrives
		if err!= nil {
			slog.Error("Accept error", "err", err)
			continue
		}
		go s.handleConn(conn)  // Handle connection in a new goroutine
	}
}


func (s *Server) handleRawMessage(rawMsg []byte) error{
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}
	
	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	}
	return nil
}

func (s *Server) handleConn(conn net.Conn){
	peer := NewPeer(conn, s.msgCh)  // Create a new Peer for the connection
	s.addPeerCh <- peer	 // Register the peer via the channel
	slog.Info("New peer conncted", "rempteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("Peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}
func main() {
	server := NewServer(Config{}) // Create server with default config
	go func() {
		log.Fatal(server.Start())  // Start the server (exit on error)
	}()
	time.Sleep(time.Second)

	for i:= 0; i<10; i++ {
		c := client.New("localhost:5001")
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(server.kv.data)

	time.Sleep(time.Second)
	//select {} //we are blocking here so the program does not exit
}