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

	sh      *test.ServerHelper
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

		sh:      sh,
		testing: t,
	}
}

// Start starts the server.
func (sh *ServerHelper) Start() *ServerHelper {
	sh.sh.Start()
	return sh
}
