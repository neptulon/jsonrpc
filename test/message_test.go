package test

import (
	"sync"
	"testing"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/jsonrpc/middleware"
)

type echoMsg struct {
	Message string `json:"message"`
}

func TestEcho(t *testing.T) {
	sh := NewServerHelper(t).Start()
	defer sh.Close()

	rout := sh.GetRouter()
	rout.Request("echo", middleware.Echo)

	ch := sh.GetClientHelper()
	defer ch.Close()

	// todo2: Helper.Middleware function should do the wg.Add(1)/wg.Done() and Close should wait for it. Also in neptulon
	var wg sync.WaitGroup
	wg.Add(1)

	ch.SendRequest("echo", echoMsg{Message: "Hello!"}, func(ctx *jsonrpc.ResCtx) error {
		defer wg.Done()
		var msg echoMsg
		if err := ctx.Result(&msg); err != nil {
			t.Fatal(err)
		}
		if msg.Message != "Hello!" {
			t.Fatalf("expected: %v got: %v", "Hello!", msg.Message)
		}
		return ctx.Next()
	})

	wg.Wait()
}

func TestOrderedDuplex(t *testing.T) {
	// client requests, server answers, server requests, client answers, client closes gracefully
}

func TestSimultaneousDuplex(t *testing.T) {
	// mash of requests, notifications, responses flying in both directions in a random loop with predefined answers
	// along with long running background request-response with large data to test interleaving and to experiment with future streaming semantics
}
