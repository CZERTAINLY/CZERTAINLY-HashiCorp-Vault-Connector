package utils

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"math/big"

	"github.com/google/uuid"
)

var log = logger.Get()

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
	uuidByte, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Error(err.Error())
	}

	return uuidByte.String()
}

func ExtractCommonName(csr string) (string, error) {
	decodedCsr, _ := base64.StdEncoding.DecodeString(csr)

	block, _ := pem.Decode(decodedCsr)
	if block == nil {
		log.Error("Failed to parse PEM block containing the CSR")
		return "", fmt.Errorf("failed to parse PEM block containing the CSR")
	}

	csrParsed, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		log.Error("Failed to parse CSR: " + err.Error())
		return "", fmt.Errorf("failed to parse CSR: %v", err)
	}

	commonName := csrParsed.Subject.CommonName
	return commonName, nil
}

func ExtractSerialNumber(certificate string) (*big.Int, error) {
	block, err := pem.Decode([]byte(certificate))
	if err != nil {
		log.Error("Failed to parse PEM block containing the certificate")
		return nil, fmt.Errorf("failed to parse PEM block containing the certificate")
	}
	if block == nil {
		log.Error("Failed to parse PEM block containing the certificate")
		return nil, fmt.Errorf("failed to parse PEM block containing the certificate")
	}

	certificateParsed, errParse := x509.ParseCertificate(block.Bytes)
	if errParse != nil {
		log.Error("Failed to parse certificate: " + errParse.Error())
	}
	serialNumber := certificateParsed.SerialNumber
	return serialNumber, nil
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
