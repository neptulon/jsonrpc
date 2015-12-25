package client

import "github.com/neptulon/neptulon/client"

// Client is a Neptulon JSON-RPC client.
type Client struct {
	client *client.Client // Inner Neptulon client.
}
