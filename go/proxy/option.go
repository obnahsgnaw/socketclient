package proxy

import (
	"github.com/obnahsgnaw/application/pkg/security"
	"github.com/obnahsgnaw/socketclient/go/auth"
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	security2 "github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
)

type Option func(*Server)

func Auth(as *auth.Auth) Option {
	return func(s *Server) {
		s.auth = as
	}
}

func Es(tp security.EsType, m security.EsMode) Option {
	return func(s *Server) {
		if tp > 0 {
			s.es = security.NewEsCrypto(tp, m)
			s.es.SetEncoder(s.encoder)
		}
	}
}

func Encoder(c security.Encoder) Option {
	return func(s *Server) {
		if c != nil {
			s.encoder = c
			s.es.SetEncoder(c)
		}
	}
}

func Encode(e bool) Option {
	return func(s *Server) {
		s.encode = e
	}
}

func JsonProxy() Option {
	return func(s *Server) {
		s.proxyDataCoder = codec.NewJsonDataBuilder()
	}
}

func ProtoProxy() Option {
	return func(s *Server) {
		s.proxyDataCoder = codec.NewProtobufDataBuilder()
	}
}

func GatewayErrHandler(f func(status gatewayv1.GatewayError_Status, triggerId uint32)) Option {
	return func(s *Server) {
		if f != nil {
			s.gatewayErrHandler = f
		}
	}
}

func Target(target *security2.Target) Option {
	return func(s *Server) {
		if s.target != nil {
			s.target = target
			if s.target.Type == "" {
				s.target.Type = "user"
			}
		}
	}
}
