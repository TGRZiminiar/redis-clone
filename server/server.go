package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"tgrziminiar/redisclone/cache"
	"tgrziminiar/redisclone/client"
	"tgrziminiar/redisclone/controller"
	"time"
)

type (
	ServerConfigs struct {
		ListenAddr string
		LeaderAddr string
		IsLeader   bool
	}

	Server struct {
		*ServerConfigs

		cache cache.Cacher

		clients map[*client.Client]struct{}
	}
)

func NewServer(cfg *ServerConfigs, cache cache.Cacher) *Server {
	return &Server{
		ServerConfigs: cfg,
		cache:         cache,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("server start failed : %s", err)
	}

	log.Printf("server start on port %s", s.ListenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error orcured on accept connection : %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {

	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error : %s\n", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))

	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *controller.CommandSet:
		s.handleSetCommand(conn, v)
		// case *proto.CommandGet:
		// s.handleGetCommand(conn, v)
		// case *proto.CommandJoin:
		// s.handleJoinCommand(conn, v)
	}
}
func (s *Server) handleSetCommand(conn net.Conn, ctrl *controller.CommandSet) error {
	log.Printf("SET %s to %s", ctrl.Key, ctrl.Value)

	go func() {
		for c := range s.clients {
			err := c.Set(context.TODO(), ctrl.Key, ctrl.Value, ctrl.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := controller.ResponseSet{}

	if err := s.cache.Set(ctrl.Key, ctrl.Value, time.Duration(ctrl.TTL)); err != nil {
		resp.Status = controller.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = controller.StatusOk
	_, err := conn.Write(resp.Bytes())

	return err
}
