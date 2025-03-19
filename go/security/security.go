package security

import (
	"errors"
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
	"github.com/obnahsgnaw/goutils/security/rsautil"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketclient/go/client"
	client2 "github.com/obnahsgnaw/socketutil/service/client"
	"go.uber.org/zap/zapcore"
)

const (
	SuccessWithoutSecurity = "000"
	SuccessWithSecurity    = "111"
	FailedWithSecurity     = "222"
)

// Server Gateway Security Control Service
type Server struct {
	client.Clienter
	rdServer    base.Server
	rsa         *rsautil.Rsa
	es          *esutil.ADes
	encoder     coder.Encoder
	encode      bool
	esKey       []byte
	initialized bool
	disabled    bool
	failedCb    func(error)
	target      *Target
}

type Target struct {
	Type    string
	Id      string
	PubCert []byte
}

func New(c client.Clienter, target *Target, o ...Option) *Server {
	s := &Server{
		Clienter: c,
		rdServer: base.Server{},
		rsa:      rsautil.New(),
		es:       esutil.New(esutil.Aes256, esutil.CbcMode),
		target:   target,
	}
	s.with(o...)
	s.withInterceptor()
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
	s.Log(zapcore.InfoLevel, "security init start")
	s.esKey = s.es.Type().RandKey()
	var encodeKey []byte
	if s.target == nil {
		s.Log(zapcore.ErrorLevel, "security invalid target")
		return
	}
	if len(s.target.PubCert) > 0 {
		var err error
		if encodeKey, err = BuildEsKeyPackage(s.rsa, s.target.PubCert, s.esKey); err != nil {
			s.Log(zapcore.ErrorLevel, "security rsa encrypt failed: "+err.Error())
			s.failedCb(errors.New("rsa encrypt failed: " + err.Error()))
		}
	}
	encodeKey = AuthenticatePackage(s.target.Type, s.target.Id, s.Config().DataCoder.Name(), encodeKey)

	if err := s.Client().SendRaw(encodeKey); err != nil {
		s.Log(zapcore.ErrorLevel, "security send initialize package failed: "+err.Error())
		s.failedCb(errors.New("send initialize package failed: " + err.Error()))
	}
}

func (s *Server) stop() {
	s.Log(zapcore.InfoLevel, "security stop")
	s.rdServer.Pause()
	s.initialized = false
	s.disabled = false
}

func (s *Server) withInterceptor() {
	s.Client().With(client2.ListenInterceptor(func(bytes []byte) []byte {
		if s.initialized {
			return bytes
		}
		bStr := string(bytes)
		if bStr == SuccessWithoutSecurity || bStr == SuccessWithSecurity {
			s.initialized = true
			s.disabled = bStr == SuccessWithoutSecurity
			s.Log(zapcore.InfoLevel, "security init success")
			s.rdServer.Ready()
			return nil
		}
		if s.failedCb != nil {
			s.Log(zapcore.ErrorLevel, "security init failed with response: "+bStr)
			s.failedCb(errors.New("init failed with response: " + bStr))
		}
		return nil
	}))
	s.Client().With(client2.GatewayPkgInterceptor(NewInterceptor(func() *esutil.ADes {
		return s.es
	}, func() []byte {
		return s.esKey
	}, func() bool {
		return s.encode
	}, func() bool {
		return s.disabled
	})))
}

// WhenReady Callback handler after the service is ready
func (s *Server) WhenReady(cb func()) {
	s.rdServer.WhenReady(cb)
}

// WhenPaused Callback handler after the service is suspended
func (s *Server) WhenPaused(cb func()) {
	s.rdServer.WhenPaused(cb)
}
