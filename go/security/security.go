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
	base.Server
	client      *client.Client
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

func New(c *client.Client, target *Target, o ...Option) *Server {
	s := &Server{
		client: c,
		rsa:    rsautil.New(),
		es:     esutil.New(esutil.Aes256, esutil.CbcMode),
		target: target,
	}
	s.With(o...)
	s.withInterceptor()
	c.WhenReady(s.start)
	c.WhenPaused(s.stop)
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
	s.client.Log(zapcore.InfoLevel, "security: init start")
	s.esKey = s.es.Type().RandKey()
	var encodeKey []byte
	if s.target == nil {
		s.client.Log(zapcore.ErrorLevel, "security: invalid target")
		return
	}
	if len(s.target.PubCert) > 0 {
		var err error
		if encodeKey, err = BuildEsKeyPackage(s.rsa, s.target.PubCert, s.esKey, s.encode); err != nil {
			s.client.Log(zapcore.ErrorLevel, "security: rsa encrypt failed: "+err.Error())
			s.failedCb(errors.New("rsa encrypt failed: " + err.Error()))
		}
	}
	encodeKey = AuthenticatePackage(s.target.Type, s.target.Id, s.client.Config().DataCoder.Name(), encodeKey)

	if err := s.client.Client().SendRaw(encodeKey); err != nil {
		s.client.Log(zapcore.ErrorLevel, "security: send initialize package failed: "+err.Error())
		s.failedCb(errors.New("send initialize package failed: " + err.Error()))
	}
}

func (s *Server) stop() {
	s.client.Log(zapcore.InfoLevel, "security: stop")
	s.Pause()
	s.initialized = false
	s.disabled = false
}

func (s *Server) withInterceptor() {
	s.client.Client().With(client2.ListenInterceptor(func(bytes []byte) []byte {
		if s.initialized {
			return bytes
		}
		bStr := string(bytes)
		if bStr == SuccessWithoutSecurity || bStr == SuccessWithSecurity {
			s.initialized = true
			s.disabled = bStr == SuccessWithoutSecurity
			s.client.Log(zapcore.InfoLevel, "security: init success")
			s.Ready()
			return nil
		}
		if s.failedCb != nil {
			s.client.Log(zapcore.ErrorLevel, "security: init failed with response: "+bStr)
			s.failedCb(errors.New("init failed with response: " + bStr))
		}
		return nil
	}))
	s.client.Client().With(client2.GatewayPkgInterceptor(NewInterceptor(func() *esutil.ADes {
		return s.es
	}, func() []byte {
		return s.esKey
	}, func() bool {
		return s.encode
	}, func() bool {
		return s.disabled
	})))
}
