// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/neptulon/neptulon"
)

// Server is a Neptulon JSON-RPC server.
type Server struct {
	neptulon *neptulon.Server
	mw       Middleware
}

// NewServer creates a Neptulon JSON-RPC server.
func NewServer(s *neptulon.Server) (*Server, error) {
	if s == nil {
		return nil, errors.New("given Neptulon server instance is nil")
	}

	rpc := Server{neptulon: s}
	s.MiddlewareIn(rpc.mw.NeptulonMiddlewareIn)
	return &rpc, nil
}

// ReqMiddleware registers middleware to handle incoming request messages.
func (s *Server) ReqMiddleware(reqMiddleware ...func(ctx *ReqCtx) error) {
	s.mw.ReqMiddleware(reqMiddleware...)
}

// NotMiddleware registers middleware to handle incoming notification messages.
func (s *Server) NotMiddleware(notMiddleware ...func(ctx *NotCtx) error) {
	s.mw.NotMiddleware(notMiddleware...)
}

// ResMiddleware registers middleware to handle incoming response messages.
func (s *Server) ResMiddleware(resMiddleware ...func(ctx *ResCtx) error) {
	s.mw.ResMiddleware(resMiddleware...)
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
