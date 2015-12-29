package jsonrpc

import (
	"encoding/json"

	"github.com/neptulon/cmap"
	"github.com/neptulon/shortid"
)

// Sender is a JSON-RPC middleware for sending requests and handling responses asynchronously.
type Sender struct {
	send      func(connID string, msg []byte) error
	resRoutes *cmap.CMap // message ID (string) -> handler func(ctx *ResCtx) error : expected responses for requests that we've sent
}

// NewSender creates a new Sender middleware.
func NewSender(send func(connID string, msg []byte) error) Sender {
	return Sender{
		send:      send,
		resRoutes: cmap.New(),
	}
}

// SendRequest sends a JSON-RPC request throught the connection denoted by the connection ID with an auto generated request ID.
// resHandler is called when a response is returned.
func (s *Sender) SendRequest(connID string, method string, params interface{}, resHandler func(ctx *ResCtx)) (reqID string, err error) {
	id, err := shortid.UUID()
	if err != nil {
		return "", err
	}

	req := Request{ID: id, Method: method, Params: params}
	if err = s.sendMsg(connID, req); err != nil {
		return "", err
	}

	s.resRoutes.Set(req.ID, resHandler)
	return id, nil
}

// SendRequestArr sends a JSON-RPC request throught the connection denoted by the connection ID, with array params and auto generated request ID.
// resHandler is called when a response is returned.
func (s *Sender) SendRequestArr(connID string, method string, resHandler func(ctx *ResCtx), params ...interface{}) (reqID string, err error) {
	return s.SendRequest(connID, method, params, resHandler)
}

// SendNotification sends a JSON-RPC notification throught the connection denoted by the connection ID with structured params object.
func (s *Sender) SendNotification(connID string, method string, params interface{}) error {
	return s.sendMsg(connID, Notification{Method: method, Params: params})
}

// SendNotificationArr sends a JSON-RPC notification message throught the connection denoted by the connection ID with array params.
func (s *Sender) SendNotificationArr(connID string, method string, params ...interface{}) error {
	return s.SendNotification(connID, method, params)
}

// SendResponse sends a JSON-RPC response throught the connection denoted by the connection ID.
func (s *Sender) SendResponse(connID string, id string, result interface{}, err *ResError) error {
	return s.sendMsg(connID, Response{ID: id, Result: result, Error: err})
}

// SendMsg sends any JSON-RPC message throught the connection denoted by the connection ID.
func (s *Sender) sendMsg(connID string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.send(connID, data)
}

// ResMiddleware is a JSON-RPC incoming response handler middleware.
func (s *Sender) resMiddleware(ctx *ResCtx) error {
	if handler, ok := s.resRoutes.GetOk(ctx.id); ok {
		err := handler.(func(ctx *ResCtx) error)(ctx)
		s.resRoutes.Delete(ctx.id)
		if err != nil {
			return err
		}
	}

	return nil
}
