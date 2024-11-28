package main

import (
	"context"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/gateway"
	"github.com/obnahsgnaw/socketutil/codec"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pub := []byte(`-----BEGIN rsa public key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuNkRe6rD0+C44ZEirIVI
PAWK7sDPN5e4LF4ztNNg8FU9b2B6VNY09lbUbXcX9cz71vsZqdzGu9TfHXgx8niu
Wr2hsIXySgREM8EdqscriL8SyOgfA3MvQs82dKEv5HAkleR9wty/OjNnJxE8uAiN
UT0KBTQTYS1+PwRBPLRghhHiQFRWz5k0uxIhColQvQucZnuQxf3xGwRzXs1r4HAF
q68jhsOVOOLfrmMb6W/Fe/znIsX9fjLL57kp+a/eHrYQ9JosWpqU6uABmuducafI
G2LowIs4xjUyEz/gRlW1gxz9owvoDrUn5vx64a+/JyJPG4qUOqG1qCY5NgJRFRjd
IwIDAQAB
-----END rsa public key-----
`)
	token := &auth.Auth{
		AppId: "",
		Token: "",
	}
	dataType := codec.Proto

	gw := gateway.Default(ctx, "127.0.0.1", 1811, dataType, pub, token)

	gw.Start()

	select {}
}
