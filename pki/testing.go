package main

import (
	"flag"
	"os"
	"time"
	"tvpn/pki/keygen"
	"crypto/x509/pkix"
	"log"
)

var (
	cn = flag.String("cn", "", "Common Name")
	name = flag.String("name", "", "Name for output files")
	ktype = flag.String("type","", "type for the output files")
)

func main() {
	flag.Parse()

	if (len(*cn) == 0 || len(*name) == 0 || len(*ktype) == 0) {
		log.Fatal("invalid input! need all of -cn, -name, and -type")
	}

	certOutFile, err := os.Create(*name + ".crt")
	if err != nil {
		log.Fatalf("Failed to create cert file: %s",err)
	}
	keyOutFile, err := os.Create(*name + ".key")
	if err != nil {
		log.Fatalf("Failed to create key file: %s",err)
	}

	var kType keygen.KeyType
	switch *ktype {
	case "ca":
		kType = keygen.CAKey
	case "server":
		kType = keygen.ServerKey
	case "client":
		kType = keygen.ClientKey
	case "all":
		kType = keygen.All
	default:
		log.Fatal("Type error: must be one of ca, server, client, or all")
	}
	subj := pkix.Name{
		Country: []string{"US"},
		Locality: []string{"Louisville"},
		Province: []string{"Kentucky"},
		CommonName: *cn,
	}

	key, cert := keygen.GenPKIKeyPair(kType,2048,subj,time.Now(),365*24*time.Hour)

	certOutFile.Write(cert.Bytes())
	certOutFile.Close()
	keyOutFile.Write(key.Bytes())
	keyOutFile.Close()
}
