package client

import (
	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
)

// Client is a Neptulon JSON-RPC client.
type Client struct {
	client           *client.Client // Inner Neptulon client.
	inReqMiddleware  []func()
	inNotMiddleware  []func()
	outReqMiddleware []func()
	outNotMiddleware []func()
}

// ConnID is a randomly generated unique client connection ID.
func (c *Client) ConnID() string {
	return c.client.ConnID()
}

// Session is a thread-safe data store for storing arbitrary data for this connection session.
func (c *Client) Session() *cmap.CMap {
	return c.client.Session()
}
