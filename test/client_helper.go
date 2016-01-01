package test

import (
	"sync"
	"testing"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon/test"
)

// ClientHelper is a Neptulon JSON-RPC Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *jsonrpc.Client

	nepCH   *test.ClientHelper // inner Neptulon ClientHelper object
	testing *testing.T
	resWG   sync.WaitGroup
}

// NewClientHelper creates a new client helper object.
func NewClientHelper(t *testing.T, addr string) *ClientHelper {
	nepCH := test.NewClientHelper(t, addr)
	c := jsonrpc.UseClient(nepCH.Client)
	nepCH.Connect()
	return &ClientHelper{
		Client:  c,
		nepCH:   nepCH,
		testing: t,
	}
}

// SendRequest sends a JSON-RPC request through the client connection with an auto generated request ID.
// resHandler is called when a response is returned.
func (ch *ClientHelper) SendRequest(method string, params interface{}, resHandler func(ctx *jsonrpc.ResCtx) error) *ClientHelper {
	ch.resWG.Add(1)
	_, err := ch.Client.SendRequest(method, params, func(ctx *jsonrpc.ResCtx) error {
		defer ch.resWG.Done()
		return resHandler(ctx)
	})

	if err != nil {
		ch.testing.Fatal("Failed to send request:", err)
	}

	ch.resWG.Wait()
	return ch
}

// Close closes a connection.
func (ch *ClientHelper) Close() {
	ch.nepCH.Close()
	ch.resWG.Wait()
}
