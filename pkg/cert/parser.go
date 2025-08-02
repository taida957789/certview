package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
)

func ParseCertificateFile(filename string) ([]*x509.Certificate, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return ParseCertificateData(data)
}

func ParseCertificateData(data []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	if isPEM(data) {
		certificates, err := parsePEMData(data)
		if err != nil {
			return nil, err
		}
		return certificates, nil
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER certificate: %v", err)
	}
	certificates = append(certificates, cert)

	return certificates, nil
}

func isPEM(data []byte) bool {
	return strings.Contains(string(data), "-----BEGIN CERTIFICATE-----")
}

func parsePEMData(data []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate
	block, rest := pem.Decode(data)

	for block != nil {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PEM certificate: %v", err)
			}
			certificates = append(certificates, cert)
		}
		block, rest = pem.Decode(rest)
	}

	if len(certificates) == 0 {
		return nil, fmt.Errorf("no valid certificates found in PEM data")
	}

	return certificates, nil
}