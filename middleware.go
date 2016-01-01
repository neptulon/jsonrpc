package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/neptulon/neptulon"
)

// Middleware is a Neptulon middleware for handling JSON-RPC protocol and relevant JSON-RPC middleware.
type Middleware struct {
	reqMiddleware []func(ctx *ReqCtx) error
	notMiddleware []func(ctx *NotCtx) error
	resMiddleware []func(ctx *ResCtx) error
}

// ReqMiddleware registers middleware to handle request messages.
func (mw *Middleware) ReqMiddleware(reqMiddleware ...func(ctx *ReqCtx) error) {
	mw.reqMiddleware = append(mw.reqMiddleware, reqMiddleware...)
}

// NotMiddleware registers middleware to handle notification messages.
func (mw *Middleware) NotMiddleware(notMiddleware ...func(ctx *NotCtx) error) {
	mw.notMiddleware = append(mw.notMiddleware, notMiddleware...)
}

// ResMiddleware registers middleware to handle response messages.
func (mw *Middleware) ResMiddleware(resMiddleware ...func(ctx *ResCtx) error) {
	mw.resMiddleware = append(mw.resMiddleware, resMiddleware...)
}

// NeptulonMiddleware handles raw messages,
// categorizes the messages as one of the three JSON-RPC message types (if they are so),
// and triggers relevant middleware.
func (mw *Middleware) neptulonMiddleware(ctx *neptulon.Ctx) error {
	var m message
	if err := json.Unmarshal(ctx.Msg, &m); err != nil {
		return fmt.Errorf("cannot deserialize message: %v", err)
	}

	// if the message is a request or response
	if m.ID != "" {
		// if the message is a request
		if m.Method != "" {
			return newReqCtx(m.ID, m.Method, m.Params, ctx.Client, mw.reqMiddleware, ctx.Session()).Next()
		}

		// if the message is a response
		return newResCtx(m.ID, m.Result, ctx.Client, mw.resMiddleware, ctx.Session()).Next()
	}

	// if the message is a notification
	if m.Method != "" {
		return newNotCtx(m.Method, m.Params, ctx.Client, mw.notMiddleware, ctx.Session()).Next()
	}

	// not a JSON-RPC message so do nothing
	return nil
}
