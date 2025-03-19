package gateway

import (
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"time"
)

type Option func(*Server)

func Error(f func(act uint32, status gatewayv1.GatewayError_Status)) Option {
	return func(s *Server) {
		s.errorCb = f
	}
}

func Heartbeat(interval time.Duration) Option {
	return func(s *Server) {
		if interval < HeartbeatMin {
			interval = HeartbeatMin
		}
		s.heartbeatInterval = interval
	}
}
