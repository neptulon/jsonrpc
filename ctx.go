package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon/client"
)

/*
 * Context object definitions for Request, Response, and Notification middleware.
 */

// ReqCtx encapsulates connection, request, and reponse objects.
type ReqCtx struct {
	Res    interface{} // Response to be returned.
	Err    *ResError   // Error to be returned.
	Client *Client     // Client connection.

	id     string          // message ID
	method string          // called method
	params json.RawMessage // request parameters

	mw      []func(ctx *ReqCtx) error
	mwIndex int
	session *cmap.CMap
}

func newReqCtx(id, method string, params json.RawMessage, client *client.Client, mw []func(ctx *ReqCtx) error, session *cmap.CMap) *ReqCtx {
	// append the last middleware to stack, which will write the response to connection, if any
	mw = append(mw, func(ctx *ReqCtx) error {
		if ctx.Res != nil || ctx.Err != nil {
			return ctx.Client.WriteResponse(ctx.id, ctx.Res, ctx.Err)
		}

		return nil
	})

	return &ReqCtx{Client: useClient(client), id: id, method: method, params: params, mw: mw}
}

// Session is a data store for storing arbitrary data within this context to communicate with other middleware handling this message.
func (ctx *ReqCtx) Session() *cmap.CMap {
	return ctx.session
}

// Params reads request parameters into given object.
// Object should be passed by reference.
func (ctx *ReqCtx) Params(v interface{}) error {
	if ctx.params != nil {
		if err := json.Unmarshal(ctx.params, v); err != nil {
			return fmt.Errorf("cannot deserialize request params: %v", err)
		}
	}

	return nil
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
	Client *Client

	method string          // called method
	params json.RawMessage // notification parameters

	mw      []func(ctx *NotCtx) error
	mwIndex int
	session *cmap.CMap
}

func newNotCtx(method string, params json.RawMessage, client *client.Client, mw []func(ctx *NotCtx) error, session *cmap.CMap) *NotCtx {
	return &NotCtx{Client: useClient(client), method: method, params: params, mw: mw}
}

// Session is a data store for storing arbitrary data within this context to communicate with other middleware handling this message.
func (ctx *NotCtx) Session() *cmap.CMap {
	return ctx.session
}

// Params reads response parameters into given object.
// Object should be passed by reference.
func (ctx *NotCtx) Params(v interface{}) error {
	if ctx.params != nil {
		if err := json.Unmarshal(ctx.params, v); err != nil {
			return fmt.Errorf("cannot deserialize notification params: %v", err)
		}
	}

	return nil
}

// Next executes the next middleware in the middleware stack.
func (ctx *NotCtx) Next() error {
	ctx.mwIndex++

	if ctx.mwIndex <= len(ctx.mw) {
		return ctx.mw[ctx.mwIndex-1](ctx)
	}

	return nil
}

// ResCtx encapsulates connection and response objects.
type ResCtx struct {
	Client *Client

	id     string          // message ID
	result json.RawMessage // result parameters

	err *resError // response error (if any)

	mw      []func(ctx *ResCtx) error
	mwIndex int
	session *cmap.CMap
}

func newResCtx(id string, result json.RawMessage, client *client.Client, mw []func(ctx *ResCtx) error, session *cmap.CMap) *ResCtx {
	return &ResCtx{Client: useClient(client), id: id, result: result, mw: mw}
}

// Session is a data store for storing arbitrary data within this context to communicate with other middleware handling this message.
func (ctx *ResCtx) Session() *cmap.CMap {
	return ctx.session
}

// Result reads response result data into given object.
// Object should be passed by reference.
func (ctx *ResCtx) Result(v interface{}) error {
	if ctx.result == nil {
		if err := json.Unmarshal(ctx.result, v); err != nil {
			return fmt.Errorf("cannot deserialize response result: %v", err)
		}
	}

	return nil
}

// Next executes the next middleware in the middleware stack.
func (ctx *ResCtx) Next() error {
	ctx.mwIndex++

	if ctx.mwIndex <= len(ctx.mw) {
		return ctx.mw[ctx.mwIndex-1](ctx)
	}

	return nil
}
