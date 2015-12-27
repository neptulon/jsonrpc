package test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon/client"
	"github.com/neptulon/neptulon/test"
)

type echoMsg struct {
	To      string `json:"to"`
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
	msg := []byte("test message")

	ch := sh.GetTCPClientHelper().MiddlewareIn(func(ctx *client.Ctx) error {
		defer wg.Done()
		if !reflect.DeepEqual(ctx.Msg, msg) {
			t.Fatalf("expected: '%s', got: '%s'", msg, ctx.Msg)
		}
		return ctx.Next()
	}).Connect()
	defer ch.Close()

	wg.Add(1)
	ch.Send(msg)
	wg.Wait()
}
