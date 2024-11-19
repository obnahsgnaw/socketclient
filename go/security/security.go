package security

import (
	"errors"
	"github.com/obnahsgnaw/application/pkg/security"
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
	rsa         *security.RsaCrypto
	es          *security.EsCrypto
	encoder     security.Encoder
	encode      bool
	publicKey   []byte
	esKey       []byte
	initialized bool
	disabled    bool
	failedCb    func(error)
	targetType  string
	targetId    string
}

func New(c *client.Client, publicKey []byte, o ...Option) *Server {
	s := &Server{
		client:     c,
		rsa:        security.NewRsa(),
		es:         security.NewEsCrypto(security.Aes256, security.CbcMode),
		publicKey:  publicKey,
		targetType: "user",
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
	s.client.Log(zapcore.InfoLevel, "security: init start")
	s.esKey = s.es.Type().RandKey()
	var encodeKey []byte
	if len(s.publicKey) > 0 {
		var err error
		if encodeKey, err = BuildEsKeyPackage(s.rsa, s.publicKey, s.esKey, s.encode); err != nil {
			s.client.Log(zapcore.ErrorLevel, "security: rsa encrypt failed: "+err.Error())
			s.failedCb(errors.New("rsa encrypt failed: " + err.Error()))
		}
	}
	encodeKey = BuildDataTypePackage(s.targetType, s.targetId, s.client.Config().DataCoder.Name(), encodeKey)

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
	s.client.Client().With(client2.GatewayPkgInterceptor(NewInterceptor(func() *security.EsCrypto {
		return s.es
	}, func() []byte {
		return s.esKey
	}, func() bool {
		return s.encode
	}, func() bool {
		return s.disabled
	})))
}
