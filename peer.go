package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msg chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msg,
		delCh: delCh,
	}
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
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			var cmd Command
			rawCmd := v.Array()[0].String()
			switch rawCmd {

			case CommandClient:
				cmd = ClientCommand{
					val: v.Array()[1].String(),
				}
			case CommandSET:
				fmt.Println("set coming in ")
				if len(v.Array()) != 3 {
					return fmt.Errorf("invalid length of variable SET command")
				}
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}

			case CommandGET:
				if len(v.Array()) != 2 {
					return fmt.Errorf("invalid length of variable SET command")
				}
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}

			case CommandHELLO:
				cmd = HelloCommand{
					val: v.Array()[1].String(),
				}
			default:
				fmt.Println("unhandle cmd coming in here: ", rawCmd)
			}
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}

		}
	}
	return nil
	// return nil, fmt.Errorf("invalid unknow command received: %s", raw)

}
