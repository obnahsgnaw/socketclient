package base

import (
	"time"
)

// Server Base listener service
type Server struct {
	readyCbs []func()
	pauseCbs []func()
}

// Ready Once the service is ready, it is called to notify other listeners to start working
func (s *Server) Ready() {
	for _, ss := range s.readyCbs {
		ss()
		time.Sleep(time.Millisecond * 200)
	}
}

// Pause When the current service is suspended, the call notifies other listeners that the service is suspended
func (s *Server) Pause() {
	for _, ss := range s.pauseCbs {
		ss()
	}
}

// WhenReady Callback handler after the service is ready
func (s *Server) WhenReady(cb func()) {
	if cb != nil {
		s.readyCbs = append(s.readyCbs, cb)
	}
}

// WhenPaused Callback handler after the service is suspended
func (s *Server) WhenPaused(cb func()) {
	if cb != nil {
		s.pauseCbs = append(s.pauseCbs, cb)
	}
}
