package action

import (
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
)

var (
	PkgErrRespAction = toAction(gatewayv1.ActionId_GatewayErr)
	PingAction       = toAction(gatewayv1.ActionId_Ping)
	PoneAction       = toAction(gatewayv1.ActionId_Pong)
	AuthReqAction    = toAction(gatewayv1.ActionId_AuthReq)
	AuthRespAction   = toAction(gatewayv1.ActionId_AuthResp)
)

func toAction(id gatewayv1.ActionId) codec.Action {
	return codec.NewAction(codec.ActionId(id), id.String())
}
