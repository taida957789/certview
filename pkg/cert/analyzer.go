package cert

import (
	"crypto/x509"
	"fmt"
	"time"
)

type CertificateInfo struct {
	Certificate    *x509.Certificate
	Subject        string
	Issuer         string
	SerialNumber   string
	NotBefore      time.Time
	NotAfter       time.Time
	IsExpired      bool
	IsCA           bool
	KeyUsage       []string
	ExtKeyUsage    []string
	SANs           []string
	SignatureAlg   string
	PublicKeyAlg   string
	PublicKeySize  int
	Extensions     []ExtensionInfo
	CrossSigns     []*x509.Certificate
}

type ExtensionInfo struct {
	OID      string
	Name     string
	Critical bool
	Value    string
}

type ChainInfo struct {
	Certificates []CertificateInfo
	IsValid      bool
	Errors       []string
	CrossSigning map[string][]*x509.Certificate
	ChainPaths   []ChainPath
}

type ChainPath struct {
	Path        []CertificateInfo
	IsComplete  bool
	Description string
}

func AnalyzeCertificateChain(certs []*x509.Certificate) *ChainInfo {
	chain := &ChainInfo{
		Certificates: make([]CertificateInfo, len(certs)),
		CrossSigning: make(map[string][]*x509.Certificate),
	}

	for i, cert := range certs {
		chain.Certificates[i] = analyzeCertificate(cert)
	}

	chain.IsValid, chain.Errors = validateChain(certs)
	chain.CrossSigning = detectCrossSigning(certs)
	chain.ChainPaths = buildChainPaths(chain)

	return chain
}

func analyzeCertificate(cert *x509.Certificate) CertificateInfo {
	info := CertificateInfo{
		Certificate:   cert,
		Subject:       cert.Subject.String(),
		Issuer:        cert.Issuer.String(),
		SerialNumber:  cert.SerialNumber.String(),
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		IsExpired:     time.Now().After(cert.NotAfter),
		IsCA:          cert.IsCA,
		SANs:          cert.DNSNames,
		SignatureAlg:  cert.SignatureAlgorithm.String(),
		PublicKeyAlg:  cert.PublicKeyAlgorithm.String(),
		Extensions:    analyzeExtensions(cert),
	}

	info.KeyUsage = parseKeyUsage(cert.KeyUsage)
	info.ExtKeyUsage = parseExtKeyUsage(cert.ExtKeyUsage)
	info.PublicKeySize = getPublicKeySize(cert)

	return info
}

func parseKeyUsage(usage x509.KeyUsage) []string {
	var usages []string
	
	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "Digital Signature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "Content Commitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "Key Encipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "Data Encipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "Key Agreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "Certificate Sign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL Sign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "Encipher Only")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "Decipher Only")
	}

	return usages
}

func parseExtKeyUsage(usage []x509.ExtKeyUsage) []string {
	var usages []string
	
	for _, u := range usage {
		switch u {
		case x509.ExtKeyUsageServerAuth:
			usages = append(usages, "Server Authentication")
		case x509.ExtKeyUsageClientAuth:
			usages = append(usages, "Client Authentication")
		case x509.ExtKeyUsageCodeSigning:
			usages = append(usages, "Code Signing")
		case x509.ExtKeyUsageEmailProtection:
			usages = append(usages, "Email Protection")
		case x509.ExtKeyUsageTimeStamping:
			usages = append(usages, "Time Stamping")
		case x509.ExtKeyUsageOCSPSigning:
			usages = append(usages, "OCSP Signing")
		default:
			usages = append(usages, fmt.Sprintf("Unknown (%d)", u))
		}
	}

	return usages
}

func getPublicKeySize(cert *x509.Certificate) int {
	switch key := cert.PublicKey.(type) {
	case interface{ Size() int }:
		return key.Size() * 8
	default:
		return 0
	}
}

func analyzeExtensions(cert *x509.Certificate) []ExtensionInfo {
	var extensions []ExtensionInfo
	
	for _, ext := range cert.Extensions {
		info := ExtensionInfo{
			OID:      ext.Id.String(),
			Critical: ext.Critical,
			Name:     getExtensionName(ext.Id.String()),
			Value:    fmt.Sprintf("%x", ext.Value),
		}
		extensions = append(extensions, info)
	}

	return extensions
}

func getExtensionName(oid string) string {
	knownOIDs := map[string]string{
		"2.5.29.14":  "Subject Key Identifier",
		"2.5.29.15":  "Key Usage",
		"2.5.29.17":  "Subject Alternative Name",
		"2.5.29.19":  "Basic Constraints",
		"2.5.29.31":  "CRL Distribution Points",
		"2.5.29.32":  "Certificate Policies",
		"2.5.29.35":  "Authority Key Identifier",
		"2.5.29.37":  "Extended Key Usage",
		"1.3.6.1.5.5.7.1.1": "Authority Information Access",
	}
	
	if name, ok := knownOIDs[oid]; ok {
		return name
	}
	return "Unknown Extension"
}

func validateChain(certs []*x509.Certificate) (bool, []string) {
	var errors []string
	
	if len(certs) == 0 {
		return false, []string{"No certificates provided"}
	}

	now := time.Now()
	for i, cert := range certs {
		if now.Before(cert.NotBefore) {
			errors = append(errors, fmt.Sprintf("Certificate %d not yet valid", i))
		}
		if now.After(cert.NotAfter) {
			errors = append(errors, fmt.Sprintf("Certificate %d has expired", i))
		}
	}

	for i := 0; i < len(certs)-1; i++ {
		child := certs[i]
		parent := certs[i+1]
		
		if err := child.CheckSignatureFrom(parent); err != nil {
			errors = append(errors, fmt.Sprintf("Certificate %d signature validation failed: %v", i, err))
		}
	}

	return len(errors) == 0, errors
}

func detectCrossSigning(certs []*x509.Certificate) map[string][]*x509.Certificate {
	crossSigns := make(map[string][]*x509.Certificate)
	
	for i, cert := range certs {
		subject := cert.Subject.String()
		issuer := cert.Issuer.String()
		
		if subject == issuer {
			continue
		}
		
		for j, otherCert := range certs {
			if i == j {
				continue
			}
			
			if areCrossSignedCertificates(cert, otherCert) {
				key := subject // Use just the subject as key
				if crossSigns[key] == nil {
					crossSigns[key] = []*x509.Certificate{cert}
				}
				
				found := false
				for _, existing := range crossSigns[key] {
					if existing == otherCert {
						found = true
						break
					}
				}
				if !found {
					crossSigns[key] = append(crossSigns[key], otherCert)
				}
			}
		}
	}
	
	filtered := make(map[string][]*x509.Certificate)
	for key, certList := range crossSigns {
		if len(certList) > 1 {
			filtered[key] = certList
		}
	}
	
	return filtered
}

func areCrossSignedCertificates(cert1, cert2 *x509.Certificate) bool {
	if cert1.Subject.String() != cert2.Subject.String() {
		return false
	}
	
	if cert1.Issuer.String() == cert2.Issuer.String() {
		return false
	}
	
	return samePublicKey(cert1, cert2)
}

func samePublicKey(cert1, cert2 *x509.Certificate) bool {
	hash1 := getPublicKeyHash(cert1)
	hash2 := getPublicKeyHash(cert2)
	
	if len(hash1) != len(hash2) {
		return false
	}
	
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	
	return true
}

func getPublicKeyHash(cert *x509.Certificate) []byte {
	switch pub := cert.PublicKey.(type) {
	case interface{ Equal(interface{}) bool }:
		return cert.RawSubjectPublicKeyInfo
	default:
		_ = pub
		return cert.RawSubjectPublicKeyInfo
	}
}

func buildChainPaths(chainInfo *ChainInfo) []ChainPath {
	var paths []ChainPath
	
	if len(chainInfo.CrossSigning) == 0 {
		// No cross-signing, build simple path
		paths = append(paths, ChainPath{
			Path:        chainInfo.Certificates,
			IsComplete:  chainInfo.IsValid,
			Description: "Primary Certificate Chain",
		})
		return paths
	}
	
	// With cross-signing, build multiple paths
	endEntity := chainInfo.Certificates[0] // Assume first cert is end entity
	
	for subject, crossSignedCerts := range chainInfo.CrossSigning {
		for i, cert := range crossSignedCerts {
			path := []CertificateInfo{endEntity}
			
			// Find the cross-signed cert in our analyzed certificates
			for _, analyzedCert := range chainInfo.Certificates {
				if analyzedCert.Subject == subject && analyzedCert.Issuer == cert.Issuer.String() {
					path = append(path, analyzedCert)
					break
				}
			}
			
			// Add root CA if it exists in our chain
			for _, analyzedCert := range chainInfo.Certificates {
				if analyzedCert.Subject == cert.Issuer.String() && analyzedCert.IsCA {
					path = append(path, analyzedCert)
					break
				}
			}
			
			description := fmt.Sprintf("Chain Path %d via %s", i+1, cert.Issuer.String())
			
			paths = append(paths, ChainPath{
				Path:        path,
				IsComplete:  len(path) > 1,
				Description: description,
			})
		}
	}
	
	return paths
}