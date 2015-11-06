# Neptulon

[![Build Status](https://travis-ci.org/neptulon/jsonrpc.svg?branch=master)](https://travis-ci.org/neptulon/jsonrpc)
[![GoDoc](https://godoc.org/github.com/neptulon/jsonrpc?status.svg)](https://godoc.org/github.com/neptulon/jsonrpc)

JSON-RPC 2.0 protocol implementation for [Neptulon](https://github.com/neptulon/neptulon) framework.

## Example

Following example creates a TLS listener with JSON-RPC 2.0 protocol and starts listening for 'ping' requests and replies with a typical 'pong'.

```go
nep, _ := neptulon.NewServer(cert, privKey, nil, "127.0.0.1:3000", true)
rpc, _ := jsonrpc.NewServer(nep)
route, _ := jsonrpc.NewRouter(rpc)

route.Request("ping", func(ctx *jsonrpc.ReqCtx) {
	ctx.Res = "pong"
})

nep.Run()
```

## Users

[Devastator](https://github.com/nbusy/devastator) mobile messaging server is written entirely using the Neptulon framework. It uses JSON-RPC 2.0 package over Neptulon to act as the server part of a mobile messaging app. You can visit its repo to see a complete use case of Neptulon framework.

## Testing

All the tests can be executed with `GORACE="halt_on_error=1" go test -race -cover ./...` command. Optionally you can add `-v` flag to observe all connection logs.

## License

[MIT](LICENSE)
