package client

import (
	"context"
	"fmt"
	"net"
	"tgrziminiar/redisclone/controller"
)

type (
	Client struct {
		conn net.Conn
	}
)

func NewClient(url string) (*Client, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &controller.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := controller.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == controller.StatusNotFound {
		return nil, fmt.Errorf("could not find key (%s)", key)
	}
	if resp.Status != controller.StatusOk {
		return nil, fmt.Errorf("server responded with non OK status [%v]", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &controller.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := controller.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}
	if resp.Status != controller.StatusOk {
		return fmt.Errorf("server responsed with non OK status [%s]", resp.Status)
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
