package jsonrpc

import (
	"encoding/json"
	"sync"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
	"github.com/neptulon/shortid"
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

// UseClient wraps an established Neptulon Client into a JSON-RPC Client.
func UseClient(client *client.Client) *Client {
	return &Client{
		client: client,
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
	c.client.MiddlewareIn(c.NeptulonMiddleware)
	return c.client.Connect(addr, debug)
}

// SendRequest writes a JSON-RPC request message to a client connection with structured params object and auto generated request ID.
func (c *Client) SendRequest(method string, params interface{}) (reqID string, err error) {
	id, err := shortid.UUID()
	if err != nil {
		return "", err
	}

	return id, c.SendMsg(Request{ID: id, Method: method, Params: params})
}

// SendRequestArr writes a JSON-RPC request message to a client connection with array params and auto generated request ID.
func (c *Client) SendRequestArr(method string, params ...interface{}) (reqID string, err error) {
	return c.SendRequest(method, params)
}

// SendNotification writes a JSON-RPC notification message to a client connection with structured params object.
func (c *Client) SendNotification(method string, params interface{}) error {
	return c.SendMsg(Notification{Method: method, Params: params})
}

// SendNotificationArr writes a JSON-RPC notification message to a client connection with array params.
func (c *Client) SendNotificationArr(method string, params ...interface{}) error {
	return c.SendNotification(method, params)
}

// SendResponse writes a JSON-RPC response message to a client connection.
func (c *Client) SendResponse(id string, result interface{}, err *ResError) error {
	return c.SendMsg(Response{ID: id, Result: result, Error: err})
}

// SendMsg writes any JSON-RPC message to a client connection.
func (c *Client) SendMsg(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.client.Send(data)
}

// Close closes a client connection.
func (c *Client) Close() error {
	return c.client.Close()
}
