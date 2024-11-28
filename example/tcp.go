package main

import (
	"context"
	security2 "github.com/obnahsgnaw/application/pkg/security"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway"
	gatewaypbv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	"github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
	"time"
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

	config := client.Default("127.0.0.1", 1811, dataType)
	conn := client.New(ctx, config)
	securityServer := security.New(conn, pub,
		security.Es(security2.Aes256, security2.CbcMode),
		security.Encoder(security2.B64Encoding),
		security.Encode(true),
		security.Failed(func(err error) {
			//
		}),
	)
	authServer := auth.New(conn, token,
		auth.Security(securityServer),
		auth.Failed(func(a *auth.Auth) {
			//
		}),
	)
	gateway.New(conn,
		gateway.Security(securityServer),
		gateway.Auth(authServer),
		gateway.Heartbeat(time.Second),
		gateway.Error(func(act uint32, status gatewaypbv1.GatewayError_Status) {
			//
		}),
	)

	conn.Start()

	select {}
}
