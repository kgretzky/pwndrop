package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/kgretzky/pwndrop/utils"
)

func GenerateTLSCertificate(common string) (*tls.Certificate, error) {
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	aYear := time.Duration(10*365*24) * time.Hour
	notAfter := notBefore.Add(aYear)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	if common == "" {
		common = utils.GenRandomString(8)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{},
			Locality:           []string{},
			Organization:       []string{},
			OrganizationalUnit: []string{},
			CommonName:         common,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &private_key.PublicKey, private_key)
	if err != nil {
		return nil, err
	}

	ret_tls := &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  private_key,
	}
	return ret_tls, nil
}

func LoadTLSCertificate(pub_path string, pkey_path string) (*tls.Certificate, error) {
	pkey, err := ioutil.ReadFile(pkey_path)
	if err != nil {
		return nil, fmt.Errorf("TLS private key not found at: %s", pkey_path)
	}
	pubkey, err := ioutil.ReadFile(pub_path)
	if err != nil {
		return nil, fmt.Errorf("TLS public key not found at: %s", pub_path)
	}
	cert, err := tls.X509KeyPair(pubkey, pkey)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
