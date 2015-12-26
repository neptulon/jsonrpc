package jsonrpc

import (
	"encoding/json"
	"log"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
)

/*
 * Context object definitions for Request, Response, and Notification middleware.
 */

// ReqCtx encapsulates connection, request, and reponse objects.
type ReqCtx struct {
	Res  interface{} // Response to be returned.
	Err  *ResError   // Error to be returned.
	Conn *Conn

	id     string          // message ID
	method string          // called method
	params json.RawMessage // request parameters

	mw      []func(ctx *ReqCtx) error
	mwIndex int
	session *cmap.CMap
}

func newReqCtx(id, method string, params json.RawMessage, client *client.Client, mw []func(ctx *ReqCtx) error) *ReqCtx {
	// append the last middleware to stack, which will write the response to connection, if any
	mw = append(mw, func(ctx *ReqCtx) error {
		if ctx.Res != nil || ctx.Err != nil {
			return ctx.Conn.WriteResponse(ctx.id, ctx.Res, ctx.Err)
		}

		return nil
	})

	return &ReqCtx{Conn: NewConn(client), id: id, method: method, params: params, mw: mw}
}

// Params reads request parameters into given object.
// Object should be passed by reference.
func (ctx *ReqCtx) Params(v interface{}) {
	if ctx.params == nil {
		return
	}

	if err := json.Unmarshal(ctx.params, v); err != nil {
		log.Fatal("Cannot deserialize request params:", err)
	}
}

// Next executes the next middleware in the middleware stack.
func (ctx *ReqCtx) Next() error {
	ctx.mwIndex++

	if ctx.mwIndex <= len(ctx.mw) {
		return ctx.mw[ctx.mwIndex-1](ctx)
	}

	return nil
}

// NotCtx encapsulates connection and notification objects.
type NotCtx struct {
	Conn *Conn
	Done bool // If set, this will prevent further middleware from handling the request

	method string          // called method
	params json.RawMessage // notification parameters

	m  []func(ctx *NotCtx) error
	mi int
}

// Params reads response parameters into given object.
// Object should be passed by reference.
func (c *NotCtx) Params(v interface{}) {
	if c.params == nil {
		return
	}

	if err := json.Unmarshal(c.params, v); err != nil {
		log.Fatal("Cannot deserialize notification params:", err)
	}
}

// ResCtx encapsulates connection and response objects.
type ResCtx struct {
	Conn *Conn
	Done bool // if set, this will prevent further middleware from handling the request

	id     string          // message ID
	result json.RawMessage // result parameters

	err *resError // response error (if any)

	m  []func(ctx *ResCtx) error
	mi int
}

// Result reads response result data into given object.
// Object should be passed by reference.
func (c *ResCtx) Result(v interface{}) {
	if c.result == nil {
		return
	}

	if err := json.Unmarshal(c.result, v); err != nil {
		log.Fatalln("Cannot deserialize response result:", err)
	}
}
