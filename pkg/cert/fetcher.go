package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

func FetchCertificatesFromDomain(domainPort string) ([]*x509.Certificate, error) {
	host, port, err := parseHostPort(domainPort)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(host, port)

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no certificates received from %s", address)
	}

	return state.PeerCertificates, nil
}

func parseHostPort(domainPort string) (host, port string, err error) {
	if strings.Contains(domainPort, ":") {
		host, port, err = net.SplitHostPort(domainPort)
		if err != nil {
			return "", "", fmt.Errorf("invalid host:port format: %v", err)
		}
	} else {
		host = domainPort
		port = "443"
	}

	if host == "" {
		return "", "", fmt.Errorf("empty hostname")
	}

	return host, port, nil
}