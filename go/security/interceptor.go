package security

import "errors"

type Interceptor struct {
	s *Server
}

func (p *Interceptor) Encode(b []byte) ([]byte, error) {
	if p.s.disabled {
		return b, nil
	}
	b1, iv, err := p.s.es.Encrypt(b, p.s.esKey, p.s.encode)
	if err != nil {
		return nil, err
	}
	return append(iv, b1...), nil
}

func (p *Interceptor) Decode(b []byte) ([]byte, error) {
	if p.s.disabled {
		return b, nil
	}

	if len(b) < p.s.es.Type().IvLen() {
		return nil, errors.New("invalid data length")
	}
	iv := b[:p.s.es.Type().IvLen()]
	b = b[p.s.es.Type().IvLen():]
	return p.s.es.Decrypt(b, p.s.esKey, iv, p.s.encode)
}
