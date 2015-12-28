package test

import (
	"testing"

	"github.com/neptulon/jsonrpc"
)

// ClientHelper is a Neptulon JSON-RPC Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *jsonrpc.Client

	testing *testing.T
}
