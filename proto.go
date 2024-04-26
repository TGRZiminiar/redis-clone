package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

const (
	CommandSET    = "set"
	CommandGET    = "get"
	CommandHELLO  = "hello"
	CommandClient = "client"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type HelloCommand struct {
	val string
}
type ClientCommand struct {
	val string
}
type GetCommand struct {
	key []byte
}

// func parseCommand(raw string) (Command, error) {
// 	rd := resp.NewReader(bytes.NewBufferString(raw))

// 	for {
// 		v, _, err := rd.ReadValue()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// fmt.Printf("Read %s\n", v.Type())

// 		var cmd Command

// 		if v.Type() == resp.Array {
// 			for _, value := range v.Array() {

// 				switch value.String() {

// 				case CommandSET:
// 					if len(v.Array()) != 3 {
// 						return nil, fmt.Errorf("invalid length of variable SET command")
// 					}
// 					cmd = SetCommand{
// 						key: v.Array()[1].Bytes(),
// 						val: v.Array()[2].Bytes(),
// 					}
// 					return cmd, nil
// 				case CommandGET:
// 					fmt.Println(v.Array()[1])
// 					if len(v.Array()) != 2 {
// 						return nil, fmt.Errorf("invalid length of variable SET command")
// 					}
// 					cmd = GetCommand{
// 						key: v.Array()[1].Bytes(),
// 					}
// 					return cmd, nil
// 				case CommandHELLO:
// 					cmd = HelloCommand{
// 						val: v.Array()[1].String(),
// 					}
// 					return cmd, nil
// 				default:
// 					return nil, fmt.Errorf("invalid unknow command received: %s", raw)
// 				}
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("invalid unknow command received: %s", raw)
// }

func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
	}
	return buf.Bytes()
}

func respWriteOK() []byte {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteString("OK")
	return buf.Bytes()
}
