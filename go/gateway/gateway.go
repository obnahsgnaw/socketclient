package gateway

import (
	"errors"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

var HeartbeatMin = 5 * time.Second // 最小 5 秒

type Server struct {
	rdServer base.Server
	client.Clienter
	heartbeatInterval time.Duration
	errorCb           func(act uint32, status gatewayv1.GatewayError_Status)
}

func New(c client.Clienter, o ...Option) *Server {
	s := &Server{
		Clienter:          c,
		rdServer:          base.Server{},
		heartbeatInterval: 10 * time.Second,
	}
	s.with(o...)
	s.withGatewayError()
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
	s.Log(zapcore.InfoLevel, "gateway start")

	if s.heartbeatInterval > 0 {
		s.Log(zapcore.InfoLevel, "gateway withed heartbeat")
		if err := s.withHeartbeat(); err != nil {
			s.Log(zapcore.ErrorLevel, "gateway heartbeat init failed, err="+err.Error())
		}
	}
	s.rdServer.Ready()
}

func (s *Server) stop() {
	s.Log(zapcore.InfoLevel, "gateway stop")
	s.rdServer.Pause()
}

func (s *Server) PauseHeartbeat() {
	s.Client().HeartbeatPause()
}

func (s *Server) ContinueHeartbeat() {
	s.Client().HeartbeatContinue()
}

func (s *Server) withHeartbeat() error {
	s.Client().Listen(action.PoneAction, func() codec.DataPtr {
		return &gatewayv1.PongResponse{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.PongResponse)
		s.Log(zapcore.DebugLevel, "gateway received pone, server time="+data.Timestamp.AsTime().String())
		return
	})
	b, err := s.Client().Pack(action.PingAction, &gatewayv1.PingRequest{})
	if err != nil {
		return errors.New("gateway pack ping package failed,err=" + err.Error())
	}
	s.Client().Heartbeat(b, s.heartbeatInterval)
	return nil
}

func (s *Server) withGatewayError() {
	s.Client().Listen(action.PkgErrRespAction, func() codec.DataPtr {
		return &gatewayv1.GatewayError{}
	}, func(rqData codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		data := rqData.(*gatewayv1.GatewayError)
		if data.Status != gatewayv1.GatewayError_None {
			s.Log(zapcore.ErrorLevel, "gateway received gateway error: triggered action="+strconv.Itoa(int(data.TriggerAction))+",err="+data.Status.String())
			if s.errorCb != nil {
				s.errorCb(data.TriggerAction, data.Status)
			}
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
