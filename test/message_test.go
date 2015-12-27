package test

import (
	"sync"
	"testing"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon/test"
)

type echoMsg struct {
	Message string `json:"message"`
}

func TestEcho(t *testing.T) {
	sh := test.NewTCPServerHelper(t).Start()
	defer sh.Close()

	js, err := jsonrpc.NewServer(sh.Server)
	if err != nil {
		t.Fatal(err)
	}

	rout, err := jsonrpc.NewRouter(&js.Middleware)
	if err != nil {
		t.Fatal(err)
	}

	rout.Request("echo", func(ctx *jsonrpc.ReqCtx) error {
		var msg echoMsg
		if err := ctx.Params(&msg); err != nil {
			t.Fatal(err)
		}
		ctx.Res = msg
		return ctx.Next()
	})

	var wg sync.WaitGroup

	ch := sh.GetTCPClientHelper().Connect()
	defer ch.Close()

	// todo: separate echo middleware into /middleware package
	// todo2: use sender.go rather than this manual handling

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
	jc.SendRequest("echo", echoMsg{Message: "Hello!"})
	wg.Wait()
}
