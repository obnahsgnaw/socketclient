package gateway

import (
	"errors"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	gatewayv1 "github.com/obnahsgnaw/socketclient/go/gateway/gen/gateway/v1"
	"github.com/obnahsgnaw/socketclient/go/security"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

var HeartbeatMin = 1 * time.Second // 最小 5 秒

type Server struct {
	base.Server
	client            *client.Client
	sec               *security.Server
	auth              *auth.Server
	heartbeatInterval time.Duration
	errorCb           func(act uint32, status gatewayv1.GatewayError_Status)
}

func New(c *client.Client, o ...Option) *Server {
	s := &Server{
		client:            c,
		heartbeatInterval: 10 * time.Second,
	}
	s.with(o...)
	s.withGatewayError()
	if s.auth != nil {
		s.auth.WhenReady(s.start)
		s.auth.WhenPaused(s.stop)
	} else if s.sec != nil {
		s.sec.WhenReady(s.start)
		s.sec.WhenPaused(s.stop)
	} else {
		s.client.WhenReady(s.start)
		s.client.WhenPaused(s.stop)
	}

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
	s.client.Log(zapcore.InfoLevel, "gateway: start")

	if s.heartbeatInterval > 0 {
		s.client.Log(zapcore.InfoLevel, "gateway: withed heartbeat")
		if err := s.withHeartbeat(); err != nil {
			s.client.Log(zapcore.ErrorLevel, "gateway: heartbeat init failed, err="+err.Error())
		}
	}
}

func (s *Server) stop() {
	s.client.Log(zapcore.InfoLevel, "gateway: stop")
	s.Pause()
}

func (s *Server) withHeartbeat() error {
	s.client.Client().Listen(action.PoneAction, func() codec.DataPtr {
		return &gatewayv1.PongResponse{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.PongResponse)
		s.client.Log(zapcore.DebugLevel, "gateway: received pone, server time="+data.Timestamp.AsTime().String())
		return
	})
	b, err := s.client.Client().Pack(action.PingAction, &gatewayv1.PingRequest{})
	if err != nil {
		return errors.New("gateway: pack ping package failed,err=" + err.Error())
	}
	s.client.Client().Heartbeat(b, s.heartbeatInterval)
	return nil
}

func (s *Server) withGatewayError() {
	s.client.Client().Listen(action.PkgErrRespAction, func() codec.DataPtr {
		return &gatewayv1.GatewayError{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.GatewayError)
		if data.Status != gatewayv1.GatewayError_None {
			s.client.Log(zapcore.ErrorLevel, "gateway: received gateway error: triggered action="+strconv.Itoa(int(data.TriggerAction))+",err="+data.Status.String())
			if s.errorCb != nil {
				s.errorCb(data.TriggerAction, data.Status)
			}
		}
		return
	})
}
