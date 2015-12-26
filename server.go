// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/client"
)

// Server is a Neptulon JSON-RPC server.
type Server struct {
	neptulon      *neptulon.Server
	reqMiddleware []func(ctx *ReqCtx) error
	notMiddleware []func(ctx *NotCtx) error
	resMiddleware []func(ctx *ResCtx) error
}

// NewServer creates a Neptulon JSON-RPC server.
func NewServer(s *neptulon.Server) (*Server, error) {
	if s == nil {
		return nil, errors.New("given Neptulon server instance is nil")
	}

	rpc := Server{neptulon: s}
	s.MiddlewareIn(rpc.neptulonMiddleware)
	return &rpc, nil
}

// ReqMiddleware registers middleware to handle incoming request messages.
func (s *Server) ReqMiddleware(reqMiddleware ...func(ctx *ReqCtx) error) {
	s.reqMiddleware = append(s.reqMiddleware, reqMiddleware...)
}

// NotMiddleware registers middleware to handle incoming notification messages.
func (s *Server) NotMiddleware(notMiddleware ...func(ctx *NotCtx) error) {
	s.notMiddleware = append(s.notMiddleware, notMiddleware...)
}

// ResMiddleware registers middleware to handle incoming response messages.
func (s *Server) ResMiddleware(resMiddleware ...func(ctx *ResCtx) error) {
	s.resMiddleware = append(s.resMiddleware, resMiddleware...)
}

// send sends a message throught the connection denoted by the connection ID.
func (s *Server) send(connID string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error while serializing JSON-RPC message: %v", err)
	}

	err = s.neptulon.Send(connID, data)
	if err != nil {
		return fmt.Errorf("error while sending JSON-RPC message: %v", err)
	}

	return nil
}

// NeptulonMiddleware deserializes incoming messages from Neptulon server and categorizes them as JSON-RPC message types, if any.
func (s *Server) neptulonMiddleware(ctx *client.Ctx) error {
	var m message
	if err := json.Unmarshal(ctx.Msg, &m); err != nil {
		return fmt.Errorf("cannot deserialize incoming message: %v", err)
	}

	// if incoming message is a request or response
	if m.ID != "" {
		// if incoming message is a request
		if m.Method != "" {
			return newReqCtx(m.ID, m.Method, m.Params, ctx.Client, s.reqMiddleware, ctx.Session()).Next()
		}

		// if incoming message is a response
		return newResCtx(m.ID, m.Result, ctx.Client, s.resMiddleware, ctx.Session()).Next()
	}

	// if incoming message is a notification
	if m.Method != "" {
		return newNotCtx(m.Method, m.Params, ctx.Client, s.notMiddleware, ctx.Session()).Next()
	}

	// not a JSON-RPC message so do nothing
	return nil
}
