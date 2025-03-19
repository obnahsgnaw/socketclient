package auth

type Option func(*Server)

func Failed(f func(*Auth)) Option {
	return func(s *Server) {
		s.failedCb = f
	}
}
