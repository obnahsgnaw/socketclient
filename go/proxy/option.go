package proxy

import (
	"github.com/obnahsgnaw/application/pkg/security"
	"github.com/obnahsgnaw/socketclient/go/auth"
)

type Option func(*Server)

func Auth(as *auth.Auth) Option {
	return func(s *Server) {
		s.auth = as
	}
}

func PublicKey(key []byte) Option {
	return func(s *Server) {
		s.publicKey = key
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
