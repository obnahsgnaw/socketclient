package main

import (
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
	auth2 "github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	proxy2 "github.com/obnahsgnaw/socketclient/go/proxy"
	"github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
	"log"
)

func main() {
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
	auth := &auth2.Auth{}
	proxy := proxy2.New("abc", "http://127.0.0.1:8028/v1/tcp-gw/proxy", codec.Json,
		proxy2.Target(&security.Target{Type: "user", PubCert: pub}),
		proxy2.Auth(auth),
		proxy2.Es(esutil.Aes256, esutil.CbcMode),
		proxy2.Encoder(coder.B64StdEncoding),
		proxy2.Encode(true),
		proxy2.GatewayErrHandler(func(status gatewayv1.GatewayError_Status, triggerId uint32) {
			log.Println("gateway error:", status, " of action ", triggerId)
		}),
	)

	resp := gatewayv1.PongResponse{}
	//proxy.DataCoder().Pack(nil)
	respAct, respData, err := proxy.SendActionPackage(action.PingAction.Id, nil)
	if err != nil {
		log.Print(err)
		return
	}
	if err = proxy.DataCoder().Unpack(respData, &resp); err != nil {
		log.Print(err)
		return
	}
	log.Print("response action:", respAct.String(), ", service time:", resp.Timestamp.AsTime().String())

	respAct, respData, err = proxy.SendActionPackage(action.PingAction.Id, nil)
	if err != nil {
		log.Print(err)
		return
	}

	if err = proxy.DataCoder().Unpack(respData, &resp); err != nil {
		log.Print(err)
		return
	}
	log.Print("response action:", respAct.String(), ", service time:", resp.Timestamp.AsTime().String())
}
