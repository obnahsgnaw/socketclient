package gateway

import (
	"github.com/obnahsgnaw/socketclient/go/auth"
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	"github.com/obnahsgnaw/socketclient/go/security"
	"time"
)

type Option func(*Server)

func Auth(as *auth.Server) Option {
	return func(s *Server) {
		s.auth = as
	}
}

func Security(ss *security.Server) Option {
	return func(s *Server) {
		s.sec = ss
	}
}

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
