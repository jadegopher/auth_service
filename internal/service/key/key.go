package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const rsaPrivateKey = "RSA PRIVATE KEY"

type IKey interface {
	String() string
	Interface() interface{}
}

type Key struct {
	privateKey *rsa.PrivateKey
}

func (k *Key) String() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  rsaPrivateKey,
			Bytes: x509.MarshalPKCS1PrivateKey(k.privateKey),
		},
	))
}

func (k *Key) Interface() interface{} {
	return k.privateKey
}
