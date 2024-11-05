package auth

import "github.com/obnahsgnaw/socketclient/go/security"

type Option func(*Server)

func Failed(f func(*Auth)) Option {
	return func(s *Server) {
		s.failedCb = f
	}
}

func Security(c *security.Server) Option {
	return func(s *Server) {
		s.sec = c
	}
}
