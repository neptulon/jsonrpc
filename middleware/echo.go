package middleware

import "github.com/neptulon/jsonrpc"

// Echo sends incoming messages back as is.
func Echo(ctx *jsonrpc.ReqCtx) error {
	var msg interface{}
	if err := ctx.Params(&msg); err != nil {
		return err
	}

	ctx.Res = msg
	return ctx.Next()
}
