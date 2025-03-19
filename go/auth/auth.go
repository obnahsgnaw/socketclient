package auth

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
)

type Auth struct {
	AppId string
	Token string
}

type Server struct {
	client.Clienter
	rdServer base.Server
	auth     *Auth
	failedCb func(*Auth)
}

func New(c client.Clienter, auth *Auth, o ...Option) *Server {
	s := &Server{rdServer: base.Server{}, Clienter: c, auth: auth}
	s.with(o...)
	s.withAuth()
	c.WhenReady(s.start)
	c.WhenPaused(s.stop)
	return s
}

func (s *Server) with(o ...Option) {
	for _, fn := range o {
		if fn != nil {
			fn(s)
		}
	}
}

func (s *Server) start() {
	s.Log(zapcore.InfoLevel, "auth init start")
	if s.auth != nil && s.auth.AppId != "" && s.auth.Token != "" {
		if err := s.Client().Send(action.AuthReqAction, &gatewayv1.AuthRequest{Token: utils.ToStr(s.auth.AppId, " ", s.auth.Token)}); err != nil {
			s.Log(zapcore.ErrorLevel, "auth init send failed, err="+err.Error())
		}
	} else {
		s.Log(zapcore.WarnLevel, "auth init ignored with empty token")
		s.rdServer.Ready()
	}
}

func (s *Server) stop() {
	s.Log(zapcore.InfoLevel, "auth stop")
	s.rdServer.Pause()
}

func (s *Server) withAuth() {
	s.Client().Listen(action.AuthRespAction, func() codec.DataPtr {
		return &gatewayv1.AuthResponse{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.AuthResponse)
		if !data.Success {
			s.Log(zapcore.ErrorLevel, "auth failed")
			if s.failedCb != nil {
				s.failedCb(s.auth)
			}
		} else {
			s.Log(zapcore.InfoLevel, "auth success")
			s.rdServer.Ready()
		}
		return
	})
}

// WhenReady Callback handler after the service is ready
func (s *Server) WhenReady(cb func()) {
	s.rdServer.WhenReady(cb)
}

// WhenPaused Callback handler after the service is suspended
func (s *Server) WhenPaused(cb func()) {
	s.rdServer.WhenPaused(cb)
}
