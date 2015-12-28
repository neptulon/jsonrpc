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
	// todo: streamline these like test.NewServerHelper(t).GetRouter().GetClientHelper() // these could wrap other helpers or directly objects?
	sh := NewServerHelper(t)
	defer sh.Close()

	rout := sh.GetRouter()

	// -----------------

	rout.Request("echo", middleware.Echo)
	sh.Start()

	var wg sync.WaitGroup

	ch := sh.nepSH.GetTCPClientHelper().Connect()
	defer ch.Close()

	// todo: separate echo middleware into /middleware package
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
	jc.SendRequest("echo", echoMsg{Message: "Hello!"})
	wg.Wait()
}

func TestOrderedDuplex(t *testing.T) {
	// client requests, server answers, server requests, client answers, client closes gracefully
}

func TestSimultaneousDuplex(t *testing.T) {
	// mash of requests, notifications, responses flying in both directions in a random loop with predefined answers
	// along with long running background request-response with large data to test interleaving and to experiment with future streaming semantics
}
