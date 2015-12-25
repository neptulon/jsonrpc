// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
		return nil, errors.New("Given Neptulon server instance is nil.")
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

	// ############### todo: pass messages through outgoing message middleware, possibly via generalizing and reusing the neptulonMiddleware
	// also good for sharing it with the Client

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Errored while serializing JSON-RPC message: %v", err)
	}

	err = s.neptulon.Send(connID, data)
	if err != nil {
		return fmt.Errorf("Errored while sending JSON-RPC message: %v", err)
	}

	return nil
}

func (s *Server) neptulonMiddleware(ctx *client.Ctx) error {
	var m message
	if err := json.Unmarshal(ctx.Msg, &m); err != nil {
		log.Fatalln("Cannot deserialize incoming message:", err)
	}

	// if incoming message is a request or response
	if m.ID != "" {
		// if incoming message is a request
		if m.Method != "" {
			rctx := ReqCtx{m: s.reqMiddleware, Conn: NewConn(ctx.Client), id: m.ID, method: m.Method, params: m.Params}
			for _, mid := range s.reqMiddleware {
				mid(&rctx)
				if rctx.Done || rctx.Res != nil || rctx.Err != nil {
					break
				}
			}

			if rctx.Res != nil || rctx.Err != nil {
				data, err := json.Marshal(Response{ID: m.ID, Result: rctx.Res, Error: rctx.Err})
				if err != nil {
					log.Fatalln("Errored while serializing JSON-RPC response:", err)
				}

				return ctx.Client.Send(data)
			}

			return nil
		}

		// if incoming message is a response
		rctx := ResCtx{Conn: NewConn(ctx.Client), id: m.ID, result: m.Result, err: m.Error}
		for _, mid := range s.resMiddleware {
			mid(&rctx)
			if rctx.Done {
				break
			}
		}

		return nil
	}

	// if incoming message is a notification
	if m.Method != "" {
		rctx := NotCtx{Conn: NewConn(ctx.Client), method: m.Method, params: m.Params}
		for _, mid := range s.notMiddleware {
			mid(&rctx)
			if rctx.Done {
				break
			}
		}

		return nil
	}

	// not a JSON-RPC message so do nothing
	return nil
}
