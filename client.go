package jsonrpc

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
)

// Client is a Neptulon JSON-RPC client.
type Client struct {
	Middleware
	client *client.Client // Inner Neptulon client.
}

// NewClient creates a new Client object.
// msgWG = (optional) sets the given *sync.WaitGroup reference to be used for counting active gorotuines that are used for handling incoming/outgoing messages.
// disconnHandler = (optional) registers a function to handle client disconnection events.
func NewClient(msgWG *sync.WaitGroup, disconnHandler func(client *client.Client)) *Client {
	return &Client{
		client: client.NewClient(msgWG, disconnHandler),
	}
}

// ConnID is a randomly generated unique client connection ID.
func (c *Client) ConnID() string {
	return c.client.ConnID()
}

// Session is a thread-safe data store for storing arbitrary data for this connection session.
func (c *Client) Session() *cmap.CMap {
	return c.client.Session()
}

// SetDeadline set the read/write deadlines for the connection, in seconds.
func (c *Client) SetDeadline(seconds int) {
	c.client.SetDeadline(seconds)
}

// UseTLS enables Transport Layer Security for the connection.
// ca = Optional CA certificate to be used for verifying the server certificate. Useful for using self-signed server certificates.
// clientCert, clientCertKey = Optional certificate/privat key pair for TLS client certificate authentication.
// All certificates/private keys are in PEM encoded X.509 format.
func (c *Client) UseTLS(ca, clientCert, clientCertKey []byte) {
	c.client.UseTLS(ca, clientCert, clientCertKey)
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string, debug bool) error {
	c.client.MiddlewareIn(c.neptulonMiddleware)
	return c.client.Connect(addr, debug)
}

// NeptulonMiddleware handles incoming messages,
// categorizes the messages as one of the three JSON-RPC message types (if they are so),
// and triggers relevant middleware.
func (c *Client) neptulonMiddleware(ctx *client.Ctx) error {
	var m message
	if err := json.Unmarshal(ctx.Msg, &m); err != nil {
		return fmt.Errorf("cannot deserialize incoming message: %v", err)
	}

	// if incoming message is a request or response
	if m.ID != "" {
		// if incoming message is a request
		if m.Method != "" {
			return newReqCtx(m.ID, m.Method, m.Params, ctx.Client, c.reqMiddleware, ctx.Session()).Next()
		}

		// if incoming message is a response
		return newResCtx(m.ID, m.Result, ctx.Client, c.resMiddleware, ctx.Session()).Next()
	}

	// if incoming message is a notification
	if m.Method != "" {
		return newNotCtx(m.Method, m.Params, ctx.Client, c.notMiddleware, ctx.Session()).Next()
	}

	// not a JSON-RPC message so do nothing
	return nil
}
