package main

import (
	"context"
	"github.com/obnahsgnaw/socketclient/go/gateway"
	"github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	target := &security.Target{
		Type:    "user",
		Id:      "xxx",
		PubCert: nil,
	}
	dataType := codec.Proto

	gw := gateway.WsDefault(ctx, "ws://127.0.0.1:29504/wss", dataType, target, nil, nil, nil)

	gw.Start()

	select {}
}
