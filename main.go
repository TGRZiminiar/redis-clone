package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/tidwall/resp"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config    *Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	delPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message
	kv        *KV
}

func NewServer(cfg *Config) *Server {

	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	log.Println("Sever running on ", s.Config.ListenAddr)

	return s.acceptLoop()

}
func (s *Server) handleMessage(msg Message) error {

	switch v := msg.cmd.(type) {
	case ClientCommand:
		if err := resp.NewWriter(msg.peer.conn).WriteString("OK"); err != nil {
			return err
		}
	case SetCommand:

		if err := s.kv.Set(v.key, v.val); err != nil {
			return err
		}
		if err := resp.NewWriter(msg.peer.conn).WriteString("OK"); err != nil {
			return err
		}

	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("no value of this key: %s \n found", v.key)
		}
		if err := resp.NewWriter(msg.peer.conn).WriteString(string(val)); err != nil {
			return err
		}

	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
			// "version": "6.2",
			// "proto":   "3",
			// "mod":     "standalone",
			// "role":    "master",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			log.Println("hello peer send error", "err: ", err)
			return err
		}

		log.Println("hello command from client")
		// log.Println("SET command want to execute key: ", v.key, " value: ", v.val)
	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				fmt.Printf("raw message error %v", err)
			}
		case peer := <-s.addPeerCh:
			fmt.Println("peer connnected remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			fmt.Println("peer disconnnected remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		case <-s.quitCh:
			return
		}
	}
}

func (s *Server) acceptLoop() error {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Panic("accept error -> ", err)
			continue
		}

		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)

	s.addPeerCh <- peer

	// log.Println("new peer connected: remoteAddr -> ", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		fmt.Println("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}

}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of the goredis server")
	flag.Parse()

	server := NewServer(&Config{
		ListenAddr: *listenAddr,
	})

	log.Fatal(server.Start())

	// blocking here so program not going to exit
	// select {}

}
