package security

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/goutils/security/rsautil"
	"github.com/obnahsgnaw/socketutil/codec"
	"strconv"
	"time"
)

func BuildEsKeyPackage(rsa *rsautil.Rsa, publicKey, esKey []byte) ([]byte, error) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	timestampKey := append(esKey, []byte(nowStr)...)
	return rsa.Encrypt(timestampKey, publicKey, true)
}

func AuthenticatePackage(typ, id string, dt codec.Name, pkg []byte) []byte {
	return append([]byte(utils.ToStr(typ, "@", id, "@", dt.String(), "::")), pkg...)
}
