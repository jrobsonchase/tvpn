package keygen

import (
	"bytes"
	"math/big"
	"encoding/pem"
	"log"
	"crypto/x509/pkix"
	"crypto/x509"
	"time"
	"crypto/rsa"
	"crypto/rand"
)

type KeyType int
const (
	CAKey KeyType = iota
	ClientKey
	ServerKey
	All
)
func GenKey(keytype KeyType,
	bits int,
	subj pkix.Name,
	start time.Time,
	dur time.Duration) (*bytes.Buffer,*bytes.Buffer) {

	var usage x509.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	var extUsage []x509.ExtKeyUsage
	var isCa bool = false

	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalf("failed to generate private key: %s",err)
	}

	switch keytype {
	case CAKey:
		usage |= x509.KeyUsageCertSign
		isCa = true
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	case ClientKey:
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	case ServerKey:
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	default:
		usage |= x509.KeyUsageCertSign
		isCa = true
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth,x509.ExtKeyUsageServerAuth}

	}

	end := start.Add(dur)
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	if end.After(endOfTime) {
		end = endOfTime
	}

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: subj,
		NotBefore: start,
		NotAfter: end,
		KeyUsage:	usage,
		ExtKeyUsage: extUsage,
		IsCA: isCa,
		BasicConstraintsValid: true,
	}

	var keyOut bytes.Buffer
	var certOut bytes.Buffer

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	keyBytes := x509.MarshalPKCS1PrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s",err)
	}

	pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE",Bytes: certBytes})
	pem.Encode(&keyOut, &pem.Block{Type: "RSA PRIVATE KEY",Bytes: keyBytes})

	return &keyOut,&certOut
}

