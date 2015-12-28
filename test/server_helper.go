package test

import (
	"testing"

	"github.com/neptulon/jsonrpc"
)

// ServerHelper is a Neptulon JSON-RPC Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	Server *jsonrpc.Server

	testing *testing.T
}
