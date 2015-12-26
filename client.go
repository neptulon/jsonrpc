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

// WriteRequest writes a JSON-RPC request message to a client connection with structured params object and auto generated request ID.
func (c *Client) WriteRequest(method string, params interface{}) (reqID string, err error) {
	id, err := shortid.UUID()
	if err != nil {
		return "", err
	}

	return id, c.WriteMsg(Request{ID: id, Method: method, Params: params})
}

// WriteRequestArr writes a JSON-RPC request message to a client connection with array params and auto generated request ID.
func (c *Client) WriteRequestArr(method string, params ...interface{}) (reqID string, err error) {
	return c.WriteRequest(method, params)
}

// WriteNotification writes a JSON-RPC notification message to a client connection with structured params object.
func (c *Client) WriteNotification(method string, params interface{}) error {
	return c.WriteMsg(Notification{Method: method, Params: params})
}

// WriteNotificationArr writes a JSON-RPC notification message to a client connection with array params.
func (c *Client) WriteNotificationArr(method string, params ...interface{}) error {
	return c.WriteNotification(method, params)
}

// WriteResponse writes a JSON-RPC response message to a client connection.
func (c *Client) WriteResponse(id string, result interface{}, err *ResError) error {
	return c.WriteMsg(Response{ID: id, Result: result, Error: err})
}

// WriteMsg writes any JSON-RPC message to a client connection.
func (c *Client) WriteMsg(msg interface{}) error {
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
