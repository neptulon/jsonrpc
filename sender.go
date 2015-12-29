package jsonrpc

import (
	"github.com/neptulon/cmap"
	"github.com/neptulon/shortid"
)

// Sender is a JSON-RPC middleware for sending requests and handling responses asynchronously.
type Sender struct {
	send      func(connID string, msg interface{}) error
	resRoutes *cmap.CMap // message ID (string) -> handler func(ctx *ResCtx) error : expected responses for requests that we've sent
}

// NewSender creates a new Sender middleware.
func NewSender() *Sender {
	return &Sender{
		resRoutes: cmap.New(),
	}
}

// SendRequest sends a JSON-RPC request throught the connection denoted by the connection ID.
// resHandler is called when a response is returned.
func (s *Sender) SendRequest(connID string, method string, params interface{}, resHandler func(ctx *ResCtx)) error {
	id, err := shortid.UUID()
	if err != nil {
		return err
	}

	req := Request{ID: id, Method: method, Params: params}
	if err = s.send(connID, req); err != nil {
		return err
	}

	s.resRoutes.Set(req.ID, resHandler)
	return nil
}

// ResMiddleware is a JSON-RPC incoming response handler middleware.
func (s *Sender) ResMiddleware(ctx *ResCtx) error {
	if handler, ok := s.resRoutes.GetOk(ctx.id); ok {
		err := handler.(func(ctx *ResCtx) error)(ctx)
		s.resRoutes.Delete(ctx.id)
		if err != nil {
			return err
		}
	}

	return nil
}
