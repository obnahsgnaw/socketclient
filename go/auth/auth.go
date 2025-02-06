package auth

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	"github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
)

type Auth struct {
	AppId string
	Token string
}

type Server struct {
	base.Server
	client   *client.Client
	sec      *security.Server
	auth     *Auth
	failedCb func(*Auth)
}

func New(c *client.Client, auth *Auth, o ...Option) *Server {
	s := &Server{client: c, auth: auth}
	s.With(o...)
	s.withAuth()
	if s.sec != nil {
		s.sec.WhenReady(s.start)
		s.sec.WhenPaused(s.stop)
	} else {
		s.client.WhenReady(s.start)
		s.client.WhenPaused(s.stop)
	}
	return s
}

func (s *Server) With(o ...Option) {
	for _, fn := range o {
		if fn != nil {
			fn(s)
		}
	}
}

func (s *Server) start() {
	s.client.Log(zapcore.InfoLevel, "auth: init start")
	if s.auth != nil && s.auth.AppId != "" && s.auth.Token != "" {
		if err := s.client.Client().Send(action.AuthReqAction, &gatewayv1.AuthRequest{Token: utils.ToStr(s.auth.AppId, " ", s.auth.Token)}); err != nil {
			s.client.Log(zapcore.ErrorLevel, "auth: init send failed, err="+err.Error())
		}
	} else {
		s.client.Log(zapcore.WarnLevel, "auth: init ignored with empty token")
		s.Ready()
	}
}

func (s *Server) stop() {
	s.client.Log(zapcore.InfoLevel, "auth: stop")
	s.Pause()
}

func (s *Server) withAuth() {
	s.client.Client().Listen(action.AuthRespAction, func() codec.DataPtr {
		return &gatewayv1.AuthResponse{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.AuthResponse)
		if !data.Success {
			s.client.Log(zapcore.ErrorLevel, "auth: failed")
			if s.failedCb != nil {
				s.failedCb(s.auth)
			}
		} else {
			s.client.Log(zapcore.InfoLevel, "auth: success")
			s.Ready()
		}
		return
	})
}
