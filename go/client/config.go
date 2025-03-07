package client

import (
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/client"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
	"time"
)

type Config struct {
	Ip                string
	Port              int
	ProtocolCoder     codec.Codec
	GatewayPkgCoder   codec.PkgBuilder
	DataCoder         codec.DataBuilder
	RetryInterval     time.Duration
	KeepaliveInterval time.Duration
	Timeout           time.Duration
	ServerLogWatcher  func(level zapcore.Level, msg string)
	PackageLogWatcher func(msgType client.MsgType, msg string, pkg []byte)
	ActionLogWatcher  func(action codec.Action, msg string)
	Network           string
	Path              string
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

func Default(ip string, port int, dataType codec.Name) *Config {
	protocolCoder, gatewayPkgCoder, dataCoder := DataTypeCoder(dataType)
	return &Config{
		Ip:                ip,
		Port:              port,
		ProtocolCoder:     protocolCoder,
		DataCoder:         dataCoder,
		GatewayPkgCoder:   gatewayPkgCoder,
		RetryInterval:     time.Second * 10,
		KeepaliveInterval: time.Second * 5,
		Timeout:           time.Second * 5,
		ServerLogWatcher: func(level zapcore.Level, msg string) {
			log.Print("server: [", PadLen(level.String(), 5), "] ", msg)
		},
		PackageLogWatcher: func(msgType client.MsgType, msg string, pkg []byte) {
			//
		},
		ActionLogWatcher: func(action codec.Action, msg string) {
			log.Println("action: ", action.String(), msg)
		},
		Network: "tcp",
	}
}

func WsDefault(ip string, port int, dataType codec.Name) *Config {
	protocolCoder := codec.NewWebsocketCodec()
	gatewayPkgCoder := codec.NewJsonPackageBuilder(ToData, ToPkg)
	dataCoder := codec.NewJsonDataBuilder()
	if dataType != codec.Json {
		dataCoder = codec.NewProtobufDataBuilder()
	}
	return &Config{
		Ip:                ip,
		Port:              port,
		ProtocolCoder:     protocolCoder,
		DataCoder:         dataCoder,
		GatewayPkgCoder:   gatewayPkgCoder,
		RetryInterval:     time.Second * 10,
		KeepaliveInterval: time.Second * 5,
		Timeout:           time.Second * 5,
		ServerLogWatcher: func(level zapcore.Level, msg string) {
			log.Print("server: [", PadLen(level.String(), 5), "] ", msg)
		},
		PackageLogWatcher: func(msgType client.MsgType, msg string, pkg []byte) {
			//
		},
		ActionLogWatcher: func(action codec.Action, msg string) {
			log.Println("action: ", action.String(), msg)
		},
		Network: "ws",
		Path:    "/wss",
	}
}

func PadLen(str string, max int) string {
	if sp := max - len(str); sp > 0 {
		str = str + strings.Repeat(" ", sp)
	}
	return str
}
