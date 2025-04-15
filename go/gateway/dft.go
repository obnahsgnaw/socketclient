package gateway

import (
	"context"
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/security"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
	"time"
)

func Default(ctx context.Context, addr string, dataType codec.Name, target *security.Target, authToken *auth.Auth, securityOptions []security.Option, authOptions []auth.Option) *Server {
	config := client.Default(addr, dataType)
	baseConn := client.New(ctx, config)
	securityConn := security.New(baseConn, target,
		append([]security.Option{
			security.Es(esutil.Aes256, esutil.CbcMode),
			security.Encoder(coder.B64StdEncoding),
			security.Encode(true),
			security.Failed(func(err error) {
				//
			}),
		}, securityOptions...)...,
	)
	authConn := auth.New(securityConn, authToken,
		append([]auth.Option{
			auth.Failed(func(a *auth.Auth) {
				//
			}),
		}, authOptions...)...,
	)
	return New(authConn,
		Heartbeat(time.Second*5),
		Error(func(act uint32, status gatewayv1.GatewayError_Status) {
			//
		}),
	)
}

func WsDefault(ctx context.Context, addr string, dataType codec.Name, target *security.Target, authToken *auth.Auth, securityOptions []security.Option, authOptions []auth.Option) *Server {
	config := client.WsDefault(addr, dataType)
	baseConn := client.New(ctx, config)
	securityConn := security.New(baseConn, target,
		append([]security.Option{
			security.Es(esutil.Aes256, esutil.CbcMode),
			security.Encoder(coder.B64StdEncoding),
			security.Encode(true),
			security.Failed(func(err error) {
				//
			}),
		}, securityOptions...)...,
	)
	authConn := auth.New(securityConn, authToken,
		append([]auth.Option{
			auth.Failed(func(a *auth.Auth) {
				//
			}),
		}, authOptions...)...,
	)
	return New(authConn,
		Heartbeat(time.Second*5),
		Error(func(act uint32, status gatewayv1.GatewayError_Status) {
			//
		}),
	)
}
