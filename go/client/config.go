package client

import (
	"github.com/obnahsgnaw/goutils/strutil"
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/client"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
	"log"
	url2 "net/url"
	"strconv"
	"time"
)

type Config struct {
	Addr              string // tcp://127.0.0.1:29508
	ip                string
	port              int
	ProtocolCoder     codec.Codec
	GatewayPkgCoder   codec.PkgBuilder
	DataCoder         codec.DataBuilder
	RetryInterval     time.Duration
	KeepaliveInterval time.Duration
	Timeout           time.Duration
	ServerLogWatcher  func(level zapcore.Level, msg string)
	PackageLogWatcher func(msgType client.MsgType, msg string, pkg []byte)
	ActionLogWatcher  func(action codec.Action, msg string)
	network           string
	path              string
}

func (c *Config) parseAddr() {
	u, err := url2.Parse(c.Addr)
	if err != nil {
		panic("invalid address: " + c.Addr)
	}
	c.network = u.Scheme
	c.ip = u.Hostname()
	c.port, _ = strconv.Atoi(u.Port())
	c.path = u.Path
}

func ToData(pkg *codec.PKG) codec.DataPtr {
	if pkg == nil {
		return &gatewayv1.GatewayPackage{}
	}

	return &gatewayv1.GatewayPackage{
		Action: pkg.Action.Val(),
		Data:   pkg.Data,
	}
}

func ToPkg(ptr codec.DataPtr) *codec.PKG {
	p := ptr.(*gatewayv1.GatewayPackage)
	return &codec.PKG{
		Action: codec.ActionId(p.Action),
		Data:   p.Data,
	}
}

func DataTypeCoder(dataType codec.Name) (codec.Codec, codec.PkgBuilder, codec.DataBuilder) {
	if dataType == codec.Json {
		return codec.NewDelimiterCodec([]byte("\\N\\B"), []byte("\\N\\B")), codec.NewJsonPackageBuilder(ToData, ToPkg), codec.NewJsonDataBuilder()
	}
	return codec.NewLengthCodec(0xAB, 1024), codec.NewProtobufPackageBuilder(ToData, ToPkg), codec.NewProtobufDataBuilder()
}

func Default(addr string, dataType codec.Name) *Config {
	protocolCoder, gatewayPkgCoder, dataCoder := DataTypeCoder(dataType)
	return &Config{
		Addr:              addr,
		ProtocolCoder:     protocolCoder,
		DataCoder:         dataCoder,
		GatewayPkgCoder:   gatewayPkgCoder,
		RetryInterval:     time.Second * 10,
		KeepaliveInterval: time.Second * 5,
		Timeout:           time.Second * 5,
		ServerLogWatcher: func(level zapcore.Level, msg string) {
			log.Print("server: [", strutil.PadLen(level.String(), 5), "] ", msg)
		},
		PackageLogWatcher: func(msgType client.MsgType, msg string, pkg []byte) {
			//
		},
		ActionLogWatcher: func(action codec.Action, msg string) {
			log.Println("action: ", action.String(), msg)
		},
	}
}

func WsDefault(addr string, dataType codec.Name) *Config {
	protocolCoder := codec.NewWebsocketCodec()
	var gatewayPkgCoder codec.PkgBuilder = codec.NewJsonPackageBuilder(ToData, ToPkg)
	dataCoder := codec.NewJsonDataBuilder()
	if dataType != codec.Json {
		dataCoder = codec.NewProtobufDataBuilder()
		gatewayPkgCoder = codec.NewProtobufPackageBuilder(ToData, ToPkg)
	}
	return &Config{
		Addr:              addr,
		ProtocolCoder:     protocolCoder,
		DataCoder:         dataCoder,
		GatewayPkgCoder:   gatewayPkgCoder,
		RetryInterval:     time.Second * 10,
		KeepaliveInterval: time.Second * 5,
		Timeout:           time.Second * 5,
		ServerLogWatcher: func(level zapcore.Level, msg string) {
			log.Print("server: [", strutil.PadLen(level.String(), 5), "] ", msg)
		},
		PackageLogWatcher: func(msgType client.MsgType, msg string, pkg []byte) {
			//
		},
		ActionLogWatcher: func(action codec.Action, msg string) {
			log.Println("action: ", action.String(), msg)
		},
	}
}
