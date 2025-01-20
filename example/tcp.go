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

	gw := gateway.Default(ctx, "127.0.0.1", 29507, dataType, target)

	gw.Start()

	select {}
}
