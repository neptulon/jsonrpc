package jsonrpc

import (
	"sync"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon"
)

// Client is a Neptulon JSON-RPC client.
type Client struct {
	Middleware
	Conn *neptulon.Conn

	sender Sender
	client *neptulon.Client // inner Neptulon client
	router *Router
}

// NewClient creates a new Client object.
// msgWG = (optional) sets the given *sync.WaitGroup reference to be used for counting active gorotuines that are used for handling incoming/outgoing messages.
// disconnHandler = (optional) registers a function to handle client disconnection events.
func NewClient(msgWG *sync.WaitGroup, disconnHandler func(client *neptulon.Client)) *Client {
	return UseClient(neptulon.NewClient(msgWG, disconnHandler))
}

// UseClient wraps an established Neptulon Client into a JSON-RPC Client.
func UseClient(client *neptulon.Client) *Client {
	c := Client{
		Conn:   client.Conn,
		client: client,
	}
	c.client.MiddlewareIn(c.Middleware.neptulonMiddleware)
	c.sender = NewSender(&c.Middleware, func(connID string, msg []byte) error { return c.client.Send(msg) })
	return &c
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
	return c.client.Connect(addr, debug)
}

// SendRequest sends a JSON-RPC request through the client connection with an auto generated request ID.
// resHandler is called when a response is returned.
func (c *Client) SendRequest(method string, params interface{}, resHandler func(ctx *ResCtx) error) (reqID string, err error) {
	return c.sender.SendRequest("", method, params, resHandler)
}

// SendRequestArr sends a JSON-RPC request through the client connection, with array params and auto generated request ID.
// resHandler is called when a response is returned.
func (c *Client) SendRequestArr(method string, resHandler func(ctx *ResCtx) error, params ...interface{}) (reqID string, err error) {
	return c.sender.SendRequestArr("", method, resHandler, params)
}

// SendNotification sends a JSON-RPC notification through the client connection with structured params object.
func (c *Client) SendNotification(method string, params interface{}) error {
	return c.sender.SendNotification("", method, params)
}

// SendNotificationArr sends a JSON-RPC notification message through the client connection with array params.
func (c *Client) SendNotificationArr(method string, params ...interface{}) error {
	return c.sender.SendNotificationArr("", method, params)
}

// SendResponse sends a JSON-RPC response through the client connection.
func (c *Client) SendResponse(id string, result interface{}, err *ResError) error {
	return c.sender.SendResponse("", id, result, err)
}

// HandleRequest regiters a handler for incoming requests.
func (c *Client) HandleRequest(route string, handler func(ctx *ReqCtx) error) {
	c.lazyRegisterRouter()
	c.router.Request(route, handler)
}

// HandleNotification regiters a handler for incoming notifications.
func (c *Client) HandleNotification(route string, handler func(ctx *NotCtx) error) {
	c.lazyRegisterRouter()
	c.router.Notification(route, handler)
}

// Close closes a client connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Router middleware needs to be registered last for other middleware to be relevant.
func (c *Client) lazyRegisterRouter() {
	if c.router == nil {
		c.router, _ = NewRouter(&c.Middleware)
	}
}
