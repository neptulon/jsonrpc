package jsonrpc

import "errors"

// Router is a JSON-RPC message routing middleware.
type Router struct {
	reqRoutes map[string]func(ctx *ReqCtx) error // method name -> handler func(ctx *ReqCtx) error
	notRoutes map[string]func(ctx *NotCtx) error // method name -> handler func(ctx *NotCtx) error
}

// NewRouter creates a JSON-RPC router instance and registers it as a Neptulon JSON-RPC middleware.
func NewRouter(m *Middleware) (*Router, error) {
	if m == nil {
		return nil, errors.New("given JSON-RPC Middleware instance is nil")
	}

	r := Router{
		reqRoutes: make(map[string]func(ctx *ReqCtx) error),
		notRoutes: make(map[string]func(ctx *NotCtx) error),
	}

	m.ReqMiddleware(r.reqMiddleware)
	m.NotMiddleware(r.notMiddleware)
	return &r, nil
}

// Request adds a new request route registry.
func (r *Router) Request(route string, handler func(ctx *ReqCtx) error) {
	r.reqRoutes[route] = handler
}

// Notification adds a new notification route registry.
func (r *Router) Notification(route string, handler func(ctx *NotCtx) error) {
	r.notRoutes[route] = handler
}

func (r *Router) reqMiddleware(ctx *ReqCtx) error {
	if handler, ok := r.reqRoutes[ctx.method]; ok {
		return handler(ctx)
	}

	return ctx.Next()
}

func (r *Router) notMiddleware(ctx *NotCtx) error {
	if handler, ok := r.notRoutes[ctx.method]; ok {
		return handler(ctx)
	}

	return ctx.Next()
}
