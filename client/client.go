package client

import (
	"sync"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
)

// Client is a Neptulon JSON-RPC client.
type Client struct {
	client           *client.Client // Inner Neptulon client.
	inReqMiddleware  []func() error
	inNotMiddleware  []func() error
	inResMiddleware  []func() error
	outReqMiddleware []func() error
	outNotMiddleware []func() error
	outResMiddleware []func() error
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

// InReqMiddleware registers middleware to handle incoming request messages.
func (c *Client) InReqMiddleware(middleware ...func() error) {
	c.inReqMiddleware = append(c.inReqMiddleware, middleware...)
}

// InNotMiddleware registers middleware to handle incoming notification messages.
func (c *Client) InNotMiddleware(middleware ...func() error) {
	c.inNotMiddleware = append(c.inNotMiddleware, middleware...)
}

// InResMiddleware registers middleware to handle incoming response messages.
func (c *Client) InResMiddleware(middleware ...func() error) {
	c.inResMiddleware = append(c.inResMiddleware, middleware...)
}

// OutReqMiddleware registers middleware to handle outgoing request messages.
func (c *Client) OutReqMiddleware(middleware ...func() error) {
	c.outReqMiddleware = append(c.outReqMiddleware, middleware...)
}

// OutNotMiddleware registers middleware to handle outgoing notification messages.
func (c *Client) OutNotMiddleware(middleware ...func() error) {
	c.outNotMiddleware = append(c.outNotMiddleware, middleware...)
}

// OutResMiddleware registers middleware to handle outgoing response messages.
func (c *Client) OutResMiddleware(middleware ...func() error) {
	c.outResMiddleware = append(c.outResMiddleware, middleware...)
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

// NeptulonMiddleware is a Neptulon middleware to handle incoming raw messages and detect the messages type.
func (c *Client) neptulonMiddleware(ctx *client.Ctx) error {
	// var m message
	// if err := json.Unmarshal(ctx.Msg, &m); err != nil {
	// 	log.Fatalln("Cannot deserialize incoming message:", err)
	// }
	//
	// // if incoming message is a request or response
	// if m.ID != "" {
	// 	// if incoming message is a request
	// 	if m.Method != "" {
	// 		rctx := ReqCtx{m: s.reqMiddleware, Conn: NewConn(ctx.Client), id: m.ID, method: m.Method, params: m.Params}
	// 		for _, mid := range s.reqMiddleware {
	// 			mid(&rctx)
	// 			if rctx.Done || rctx.Res != nil || rctx.Err != nil {
	// 				break
	// 			}
	// 		}
	//
	// 		if rctx.Res != nil || rctx.Err != nil {
	// 			data, err := json.Marshal(Response{ID: m.ID, Result: rctx.Res, Error: rctx.Err})
	// 			if err != nil {
	// 				log.Fatalln("Errored while serializing JSON-RPC response:", err)
	// 			}
	//
	// 			return ctx.Client.Send(data)
	// 		}
	//
	// 		return nil
	// 	}
	//
	// 	// if incoming message is a response
	// 	rctx := ResCtx{Conn: NewConn(ctx.Client), id: m.ID, result: m.Result, err: m.Error}
	// 	for _, mid := range s.resMiddleware {
	// 		mid(&rctx)
	// 		if rctx.Done {
	// 			break
	// 		}
	// 	}
	//
	// 	return nil
	// }
	//
	// // if incoming message is a notification
	// if m.Method != "" {
	// 	rctx := NotCtx{Conn: NewConn(ctx.Client), method: m.Method, params: m.Params}
	// 	for _, mid := range s.notMiddleware {
	// 		mid(&rctx)
	// 		if rctx.Done {
	// 			break
	// 		}
	// 	}
	//
	// 	return nil
	// }

	// not a JSON-RPC message so do nothing
	return nil
}
