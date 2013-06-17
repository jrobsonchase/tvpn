package pki

import (
	"os"
	"time"
	"crypto/x509/pkix"
	"log"
)

func GenKeys(outDir,name,country,st,loc string) {

	certOutFile, err := os.Create(outDir + name + ".crt")
	if err != nil {
		log.Fatalf("Failed to create cert file: %s",err)
	}
	keyOutFile, err := os.Create(outDir + name + ".key")
	if err != nil {
		log.Fatalf("Failed to create key file: %s",err)
	}

	subj := pkix.Name{
		Country: []string{country},
		Locality: []string{loc},
		Province: []string{st},
		CommonName: name,
	}

	key, cert := GenPKIKeyPair(All,2048,subj,time.Now(),3650*24*time.Hour)

	certOutFile.Write(cert.Bytes())
	certOutFile.Close()
	keyOutFile.Write(key.Bytes())
	keyOutFile.Close()
}
