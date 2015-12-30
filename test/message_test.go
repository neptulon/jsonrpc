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

	// todo5: implement sh.GetClientHelper() instead after deciding what and how to wrap Neptulon Client/ClientHelper
	ch := sh.nepSH.GetTCPClientHelper().Connect()
	defer ch.Close()

	rout := sh.GetRouter()
	rout.Request("echo", middleware.Echo)

	// -----------------

	var wg sync.WaitGroup

	// todo2: use sender.go rather than this manual handling
	// todo3: Helper.Middleware function should do the wg.Add(1)/wg.Done() and Close should wait for it. Also in neptulon

	jc := jsonrpc.UseClient(ch.Client)
	jc.ResMiddleware(func(ctx *jsonrpc.ResCtx) error {
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

	wg.Add(1)
	// todo4 (do after todo3): SendRequest should use Sender automatically
	//  and accept response callback which would be registered as last middleware
	//  same goes for Server.SendTo() also, which would use the same Sender.go middleware
	jc.SendRequest("echo", echoMsg{Message: "Hello!"}, nil)
	wg.Wait()
}

func TestOrderedDuplex(t *testing.T) {
	// client requests, server answers, server requests, client answers, client closes gracefully
}

func TestSimultaneousDuplex(t *testing.T) {
	// mash of requests, notifications, responses flying in both directions in a random loop with predefined answers
	// along with long running background request-response with large data to test interleaving and to experiment with future streaming semantics
}
