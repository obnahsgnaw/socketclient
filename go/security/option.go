package security

import (
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
)

type Option func(*Server)

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

func Failed(f func(error)) Option {
	return func(s *Server) {
		if f != nil {
			s.failedCb = f
		}
	}
}
