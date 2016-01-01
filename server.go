// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"errors"

	"github.com/neptulon/neptulon"
)

// Server is a Neptulon JSON-RPC server.
type Server struct {
	Middleware
	Sender
	neptulon *neptulon.Server
}

// NewServer creates a Neptulon JSON-RPC server.
func NewServer(n *neptulon.Server) (*Server, error) {
	if n == nil {
		return nil, errors.New("given Neptulon server instance is nil")
	}

	s := Server{neptulon: n}
	n.MiddlewareIn(s.Middleware.neptulonMiddleware)
	s.Sender = NewSender(&s.Middleware, n.Send)

	return &s, nil
}
