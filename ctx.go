package jsonrpc

import (
	"encoding/json"
	"log"
)

/*
 * Context object definitions for Request, Response, and Notification middleware.
 */

// ReqCtx encapsulates connection, request, and reponse objects.
type ReqCtx struct {
	Conn *Conn
	Res  interface{} // Response to be returned
	Err  *ResError   // Error to be returned
	Done bool        // If set, this will prevent further middleware from handling the request

	id     string          // message ID
	method string          // called method
	params json.RawMessage // request parameters

	m  []func(ctx *ReqCtx)
	mi int
}

// Params reads request parameters into given object.
// Object should be passed by reference.
func (c *ReqCtx) Params(v interface{}) {
	if c.params == nil {
		return
	}

	if err := json.Unmarshal(c.params, v); err != nil {
		log.Fatal("Cannot deserialize request params:", err)
	}
}

// Next executes the next middleware in the middleware stack.
func (c *ReqCtx) Next() {
	c.mi++

	if c.mi <= len(c.m) {
		c.m[c.mi-1](c)
	} else if c.Res != nil {
		if err := c.Conn.Write(c.Res); err != nil {
			log.Fatalln("Errored while writing response to connection:", err)
		}
	}
}

// NotCtx encapsulates connection and notification objects.
type NotCtx struct {
	Conn *Conn
	Done bool // If set, this will prevent further middleware from handling the request

	method string          // called method
	params json.RawMessage // notification parameters

	m  []func(ctx *NotCtx)
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

	m  []func(ctx *ResCtx)
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
