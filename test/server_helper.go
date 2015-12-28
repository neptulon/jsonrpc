package test

import (
	"testing"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon/test"
)

// ServerHelper is a Neptulon JSON-RPC Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	Server *jsonrpc.Server

	nepSH   *test.ServerHelper // inner Neptulon ServerHelper object
	testing *testing.T
}

// NewServerHelper creates a new server helper object.
func NewServerHelper(t *testing.T) *ServerHelper {
	sh := test.NewTCPServerHelper(t)
	js, err := jsonrpc.NewServer(sh.Server)
	if err != nil {
		t.Fatal(err)
	}

	return &ServerHelper{
		Server: js,

		nepSH:   sh,
		testing: t,
	}
}

// GetRouter creates and attaches a new Router middleware to JSON-RPC server and returns it.
func (sh *ServerHelper) GetRouter() *jsonrpc.Router {
	route, err := jsonrpc.NewRouter(&sh.Server.Middleware)
	if err != nil {
		sh.testing.Fatal(err)
	}

	return route
}

// Start starts the server.
func (sh *ServerHelper) Start() *ServerHelper {
	sh.nepSH.Start()
	return sh
}

// Close stops the server listener and connections.
func (sh *ServerHelper) Close() {
	sh.nepSH.Close()
}
