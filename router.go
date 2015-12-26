package jsonrpc

import "errors"

// Router is a JSON-RPC message routing middleware.
type Router struct {
	server    *Server
	reqRoutes map[string]func(ctx *ReqCtx) error // method name -> handler func(ctx *ReqCtx) error
	notRoutes map[string]func(ctx *NotCtx) error // method name -> handler func(ctx *NotCtx) error
}

// NewRouter creates a JSON-RPC router instance and registers it with the Neptulon JSON-RPC server.
func NewRouter(s *Server) (*Router, error) {
	if s == nil {
		return nil, errors.New("given Neptulon server instance is nil")
	}

	r := Router{
		server:    s,
		reqRoutes: make(map[string]func(ctx *ReqCtx) error),
		notRoutes: make(map[string]func(ctx *NotCtx) error),
	}

	s.ReqMiddleware(r.reqMiddleware)
	s.NotMiddleware(r.notMiddleware)
	return &r, nil
}

// Request adds a new incoming request route registry.
func (r *Router) Request(route string, handler func(ctx *ReqCtx) error) {
	r.reqRoutes[route] = handler
}

// Notification adds a new incoming notification route registry.
func (r *Router) Notification(route string, handler func(ctx *NotCtx) error) {
	r.notRoutes[route] = handler
}

func (r *Router) reqMiddleware(ctx *ReqCtx) error {
	if handler, ok := r.reqRoutes[ctx.method]; ok {
		return handler(ctx)
	}

	return nil
}

func (r *Router) notMiddleware(ctx *NotCtx) error {
	if handler, ok := r.notRoutes[ctx.method]; ok {
		return handler(ctx)
	}

	return nil
}
