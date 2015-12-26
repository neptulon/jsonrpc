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
	Middleware
	neptulon *neptulon.Server
}

// NewServer creates a Neptulon JSON-RPC server.
func NewServer(s *neptulon.Server) (*Server, error) {
	if s == nil {
		return nil, errors.New("given Neptulon server instance is nil")
	}

	rpc := Server{neptulon: s}
	s.MiddlewareIn(rpc.NeptulonMiddleware)
	return &rpc, nil
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
