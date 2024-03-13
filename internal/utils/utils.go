package utils

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"github.com/google/uuid"
)

func DeterministicGUID(parts ...string) string {
	// concatenate all strings
	var combined string
	for _, part := range parts {
		combined += part
	}

	md5hash := md5.New()
	md5hash.Write([]byte(combined))

	// convert the hash value to a string
	md5string := hex.EncodeToString(md5hash.Sum(nil))

	// generate the UUID from the
	// first 16 bytes of the MD5 hash
	uuid, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Fatal(err)
	}

	return uuid.String()
}

func ExtractCommonName(csr string) string {
	decodedCsr, _ := base64.StdEncoding.DecodeString(csr)

	block, _ := pem.Decode([]byte(decodedCsr))
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing the CSR")
	}

	csrParsed, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse CSR: %v", err)
	}

	commonName := csrParsed.Subject.CommonName
	return commonName
}

func ExtractSerialNumber(certificate string) *big.Int {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing the certificate")
	}

	certificateParsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}
	serialNumber := certificateParsed.SerialNumber
	return serialNumber
}

func GetCertificatesFromDer(pemData []byte) ([]string, error) {

	var certs []string
	for {
		block, rest := pem.Decode(pemData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("failed to decode PEM block containing certificate")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}
		certBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})

		certs = append(certs, string(certBlock))
		pemData = rest
	}

	return certs, nil
}
