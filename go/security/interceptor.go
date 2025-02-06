package security

import (
	"errors"
	"github.com/obnahsgnaw/goutils/security/esutil"
)

type Interceptor struct {
	disabled func() bool
	es       func() *esutil.ADes
	esKey    func() []byte
	encode   func() bool
}

func NewInterceptor(es func() *esutil.ADes, esKey func() []byte, encode func() bool, disabled func() bool) *Interceptor {
	return &Interceptor{
		disabled: disabled,
		es:       es,
		esKey:    esKey,
		encode:   encode,
	}
}

func (p *Interceptor) Encode(b []byte) ([]byte, error) {
	if p.disabled() {
		return b, nil
	}
	b1, iv, err := p.es().Encrypt(b, p.esKey(), p.encode())
	if err != nil {
		return nil, err
	}
	return append(iv, b1...), nil
}

func (p *Interceptor) Decode(b []byte) ([]byte, error) {
	if p.disabled() {
		return b, nil
	}

	if len(b) < p.es().Type().IvLen() {
		return nil, errors.New("invalid data length")
	}
	iv := b[:p.es().Type().IvLen()]
	b = b[p.es().Type().IvLen():]
	return p.es().Decrypt(b, p.esKey(), iv, p.encode())
}
