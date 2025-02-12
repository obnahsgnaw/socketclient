package proxy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/goutils/security/coder"
	"github.com/obnahsgnaw/goutils/security/esutil"
	"github.com/obnahsgnaw/goutils/security/rsautil"
	"github.com/obnahsgnaw/socketclient/go/auth"
	"github.com/obnahsgnaw/socketclient/go/client"
	"github.com/obnahsgnaw/socketclient/go/gateway/action"
	proxyv1 "github.com/obnahsgnaw/socketclient/go/proxy/gen/tcpgw_frontend_api/proxy/v1"
	security2 "github.com/obnahsgnaw/socketclient/go/security"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
	"io"
	"log"
	"net/http"
	"time"
)

type Server struct {
	client            *http.Client
	dataType          codec.Name
	auth              *auth.Auth
	rsa               *rsautil.Rsa
	es                *esutil.ADes
	encoder           coder.Encoder
	gatewayPkgCoder   codec.PkgBuilder
	dataCoder         codec.DataBuilder
	proxyDataCoder    codec.DataBuilder
	interceptor       *security2.Interceptor
	encode            bool
	initialized       bool
	securityDisabled  bool
	esKey             []byte
	clientId          string
	proxyUrl          string
	gatewayErrHandler func(status gatewayv1.GatewayError_Status, triggerId uint32)
	target            *security2.Target
}

func New(clientId, proxyUrl string, dataType codec.Name, o ...Option) *Server {
	_, gatewayPkgCoder, dataCoder := client.DataTypeCoder(dataType)
	s := &Server{
		client: &http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 30 * time.Second,
		},
		dataType:        dataType,
		rsa:             rsautil.New(),
		es:              esutil.New(esutil.Aes256, esutil.CbcMode),
		encoder:         coder.B64StdEncoding,
		gatewayPkgCoder: gatewayPkgCoder,
		dataCoder:       dataCoder,
		proxyDataCoder:  codec.NewProtobufDataBuilder(),
		clientId:        toMd5(clientId),
		proxyUrl:        proxyUrl,
	}
	s.with(o...)
	if s.target == nil {
		s.target = &security2.Target{Type: "user"}
	}
	s.interceptor = security2.NewInterceptor(
		func() *esutil.ADes {
			return s.es
		},
		func() []byte {
			return s.esKey
		},
		func() bool {
			return s.encode
		},
		func() bool {
			return s.securityDisabled
		},
	)
	s.gatewayErrHandler = func(status gatewayv1.GatewayError_Status, triggerId uint32) {
		log.Println("Gateway Error:", status, " of trigger action ", triggerId)
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

// SendActionPackage SendRedirect Direct transparent transmission sends packets
func (s *Server) SendActionPackage(act codec.ActionId, data []byte) (codec.ActionId, []byte, error) {
	if err := s.init(); err != nil {
		return codec.ActionId(0), nil, err
	}
	respAct, respData, err := s.sendActionPackage(act, data)
	if err != nil && !s.initialized { // try again
		if err = s.init(); err != nil {
			return codec.ActionId(0), nil, err
		}
		respAct, respData, err = s.sendActionPackage(act, data)
	}
	return respAct, respData, err
}

func (s *Server) DataCoder() codec.DataBuilder {
	return s.dataCoder
}

func (s *Server) sendActionPackage(act codec.ActionId, data []byte) (respAct codec.ActionId, respData []byte, err error) {
	if data, err = s.gatewayPkgCoder.Pack(&codec.PKG{Action: act, Data: data}); err != nil {
		err = errors.New("encode gateway package failed, err=" + err.Error())
		return
	}

	if data, err = s.interceptor.Encode(data); err != nil {
		err = errors.New("encrypt data failed, err=" + err.Error())
		return
	}

	if data, err = s.request("POST", s.proxyUrl, data, false); err != nil {
		err = errors.New("send package failed, err=" + err.Error())
		return
	}

	if string(data) == security2.FailedWithSecurity {
		s.initialized = false
		err = errors.New("need init")
		return
	}

	if data, err = s.interceptor.Decode(data); err != nil {
		err = errors.New("decrypt data failed, err=" + err.Error())
		return
	}

	var pkg *codec.PKG
	if pkg, err = s.gatewayPkgCoder.Unpack(data); err != nil {
		err = errors.New("decode gateway package failed, err=" + err.Error())
		return
	}

	// gateway error
	if gatewayv1.ActionId(pkg.Action.Val()) == gatewayv1.ActionId_GatewayErr {
		err = errors.New("gateway package error")
		gwErr := gatewayv1.GatewayError{}
		if err1 := s.dataCoder.Unpack(pkg.Data, &gwErr); err1 == nil {
			s.gatewayErrHandler(gwErr.Status, gwErr.TriggerAction)
		}
		return
	}

	return pkg.Action, pkg.Data, nil
}

// Initialize data type, encryption and decryption, authentication, etc
func (s *Server) init() error {
	if !s.initialized {
		if err := s.authenticate(); err != nil {
			return err
		}
		if err := s.doAuth(); err != nil {
			return err
		}
	}
	return nil
}

// Exchange encryption and decryption keys,
func (s *Server) authenticate() (err error) {
	s.esKey = s.es.Type().RandKey()
	var pkg []byte
	if len(s.target.PubCert) > 0 {
		if pkg, err = security2.BuildEsKeyPackage(s.rsa, s.target.PubCert, s.esKey); err != nil {
			return err
		}
	}
	var resp []byte
	if resp, err = s.request("POST", s.proxyUrl, security2.AuthenticatePackage(s.target.Type, s.target.Id, s.dataType, pkg), true); err != nil {
		return err
	}
	respStatus := string(resp)
	if respStatus == security2.SuccessWithSecurity || respStatus == security2.SuccessWithoutSecurity {
		s.initialized = true
		s.securityDisabled = respStatus == security2.SuccessWithoutSecurity
		return nil
	}
	return errors.New("authenticate failed with: " + respStatus)
}

// Perform login authentication
func (s *Server) doAuth() (err error) {
	if s.auth != nil && s.auth.AppId != "" && s.auth.Token != "" {
		var data []byte
		if data, err = s.dataCoder.Pack(&gatewayv1.AuthRequest{Token: utils.ToStr(s.auth.AppId, " ", s.auth.Token)}); err != nil {
			return err
		}
		respAct, respData, err1 := s.sendActionPackage(action.AuthReqAction.Id, data)
		if err1 != nil {
			return err1
		}
		if respAct != action.AuthRespAction.Id {
			return errors.New("auth fail with response " + respAct.String())
		}

		var authResp gatewayv1.AuthResponse
		if err = s.dataCoder.Unpack(respData, &authResp); err != nil {
			return errors.New("decode auth response failed, err=" + err.Error())
		}
		if !authResp.Success {
			return errors.New("auth fail")
		}
	}
	return nil
}

func (s *Server) request(method string, url string, body []byte, init bool) (pkg []byte, err error) {
	var resp *http.Response
	if body, err = s.toProxyPackage(body, init); err != nil {
		err = errors.New("encode proxy package failed, err=" + err.Error())
		return nil, err
	}
	var forwardReq *http.Request
	forwardReq, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	dataType := "application/json"
	if s.proxyDataCoder.Name() == codec.Proto {
		//dataType = "application/x-protobuf"
		dataType = "application/octet-stream"
	}
	forwardReq.Header.Set("Content-Type", dataType)
	if resp, err = s.client.Do(forwardReq); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, _ = io.ReadAll(resp.Body)
		defer func(b io.ReadCloser) { _ = b.Close() }(resp.Body)
		err = errors.New("request failed with " + resp.Status)
		return
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body failed, err=" + err.Error())
		return
	}
	defer func(b io.ReadCloser) { _ = b.Close() }(resp.Body)
	pkg, err = s.parseProxyPackage(body)
	return
}

func (s *Server) toProxyPackage(body []byte, init bool) ([]byte, error) {
	data := &proxyv1.SendRequest{
		CodecType: 0,
		ClientId:  s.clientId,
		Package:   body,
		Init:      init,
	}
	if s.dataType == codec.Json {
		data.CodecType = proxyv1.CodecType_Json
	} else {
		data.CodecType = proxyv1.CodecType_Proto
	}

	return s.proxyDataCoder.Pack(data)
}

func (s *Server) parseProxyPackage(body []byte) ([]byte, error) {
	pkg := proxyv1.SendResponse{}
	if err := s.proxyDataCoder.Unpack(body, &pkg); err != nil {
		return nil, errors.New("decode proxy response data failed, err=" + err.Error())
	}
	return pkg.Package, nil
}

func toMd5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
