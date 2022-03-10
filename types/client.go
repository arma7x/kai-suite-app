package types

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	id string
	device string
	verify bool
	conn *websocket.Conn
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) SetId(id string) string {
	c.id = id
  return c.id
}

func (c *Client) GetDevice() string {
	return c.device
}

func (c *Client) SetDevice(device string) string {
  c.device = device
	return c.device
}

func (c *Client) GetVerify() bool {
	return c.verify
}

func (c *Client) SetVerify(verify bool) bool {
  c.verify = verify
	return c.verify
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

func CreateClient(id, device string, verify bool, conn *websocket.Conn) *Client {
	return &Client{id, device, verify, conn}
}
