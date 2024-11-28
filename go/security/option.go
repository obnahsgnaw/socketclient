package security

import (
	"github.com/obnahsgnaw/application/pkg/security"
)

type Option func(*Server)

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

func TargetInfo(target *Target) Option {
	return func(s *Server) {
		if s.target != nil {
			s.target = target
			if s.target.Type == "" {
				s.target.Type = "user"
			}
		}
	}
}

func Failed(f func(error)) Option {
	return func(s *Server) {
		if f != nil {
			s.failedCb = f
		}
	}
}
