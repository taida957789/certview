package cmd

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"certview/pkg/cert"
	"certview/pkg/html"
)

func RunCLI(input string) {
	var certs []*x509.Certificate
	var err error
	var title string

	if strings.Contains(input, ":") || (!strings.Contains(input, ".") && !strings.HasSuffix(input, ".pem") && !strings.HasSuffix(input, ".crt") && !strings.HasSuffix(input, ".cer")) {
		fmt.Fprintf(os.Stderr, "Fetching certificates from domain: %s\n", input)
		certs, err = cert.FetchCertificatesFromDomain(input)
		title = fmt.Sprintf("Domain: %s", input)
	} else {
		fmt.Fprintf(os.Stderr, "Parsing certificate file: %s\n", input)
		certs, err = cert.ParseCertificateFile(input)
		title = fmt.Sprintf("File: %s", input)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Found %d certificate(s)\n", len(certs))

	chainInfo := cert.AnalyzeCertificateChain(certs)

	htmlOutput, err := html.GenerateHTML(chainInfo, title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HTML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(htmlOutput)
}