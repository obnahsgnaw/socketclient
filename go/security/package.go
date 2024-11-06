package security

import (
	"github.com/obnahsgnaw/application/pkg/security"
	"github.com/obnahsgnaw/socketutil/codec"
	"strconv"
	"time"
)

func BuildEsKeyPackage(rsa *security.RsaCrypto, publicKey, esKey []byte, encode bool) ([]byte, error) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	timestampKey := append(esKey, []byte(nowStr)...)
	return rsa.Encrypt(timestampKey, publicKey, encode)
}

func BuildDataTypePackage(dt codec.Name, pkg []byte) []byte {
	// proto 增加协议字节
	if dt == codec.Proto && (len(pkg) == 0 || pkg[0] == 'j') {
		pkg = append([]byte("b"), pkg...)
	}
	// json 增加协议字节
	if dt == codec.Json && (len(pkg) == 0 || (pkg[0] != 'j' && pkg[0] != '{')) {
		pkg = append([]byte("j"), pkg...)
	}
	return pkg
}
