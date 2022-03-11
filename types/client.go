package types

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	device string
	conn *websocket.Conn
}

func (c *Client) GetDevice() string {
	return c.device
}

func (c *Client) SetDevice(device string) string {
  c.device = device
	return c.device
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

func CreateClient(device string, conn *websocket.Conn) *Client {
	return &Client{device, conn}
}
