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
		Type:    "device",
		Id:      "xxx",
		PubCert: nil,
	}
	dataType := codec.Proto

	gw := gateway.Default(ctx, "tcp://127.0.0.1:29508", dataType, target, nil, nil, nil)

	gw.Start()

	select {}
}
