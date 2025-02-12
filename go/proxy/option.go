package proxy

import (
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/security"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
)

type Option func(*Server)

func Auth(as *auth.Auth) Option {
	return func(s *Server) {
		s.auth = as
	}
}

func Es(tp esutil.EsType, m esutil.EsMode) Option {
	return func(s *Server) {
		if tp > 0 {
			s.es = esutil.New(tp, m, esutil.Encoder(s.encoder))
		}
	}
}

func Encoder(c coder.Encoder) Option {
	return func(s *Server) {
		if c != nil {
			s.encoder = c
			s.es = esutil.New(s.es.Type(), s.es.Mode(), esutil.Encoder(s.encoder))
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

func Target(target *security.Target) Option {
	return func(s *Server) {
		if s.target != nil {
			s.target = target
			if s.target.Type == "" {
				s.target.Type = "user"
			}
		}
	}
}
