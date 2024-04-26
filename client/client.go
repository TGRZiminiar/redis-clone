package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Addr string
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, val any) error {
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.AnyValue(val),
	})

	_, err := c.conn.Write(buf.Bytes())
	// _, err := io.Copy(c.conn, buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})

	_, err := c.conn.Write(buf.Bytes())
	// _, err := io.Copy(c.conn, buf)
	if err != nil {
		return "", err
	}

	bb := make([]byte, 1024)
	n, err := c.conn.Read(bb)
	if err != nil {
		return "", err
	}
	return string(bb[:n]), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
