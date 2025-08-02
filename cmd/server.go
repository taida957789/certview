package cmd

import (
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"certview/pkg/cert"
	"certview/pkg/html"
)

func RunServer(port int) {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/analyze", handleAnalyze)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting CertView server on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		formHTML, err := html.GenerateWebForm()
		if err != nil {
			http.Error(w, "Error generating form", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(formHTML))
	
	case http.MethodPost:
		// Handle form submission from home page
		handleAnalyze(w, r)
	
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	source := r.FormValue("source")
	var certs []*x509.Certificate
	var title string

	switch source {
	case "domain":
		domain := strings.TrimSpace(r.FormValue("domain"))
		if domain == "" {
			http.Error(w, "Domain is required", http.StatusBadRequest)
			return
		}

		certs, err = cert.FetchCertificatesFromDomain(domain)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching certificates: %v", err), http.StatusInternalServerError)
			return
		}
		title = fmt.Sprintf("Domain: %s", domain)

	case "file":
		file, header, err := r.FormFile("certfile")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file content", http.StatusInternalServerError)
			return
		}

		certs, err = cert.ParseCertificateData(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing certificate file: %v", err), http.StatusBadRequest)
			return
		}
		title = fmt.Sprintf("File: %s", header.Filename)

	case "paste":
		certData := strings.TrimSpace(r.FormValue("certdata"))
		if certData == "" {
			http.Error(w, "Certificate data is required", http.StatusBadRequest)
			return
		}

		certs, err = cert.ParseCertificateData([]byte(certData))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing certificate data: %v", err), http.StatusBadRequest)
			return
		}
		title = "Pasted Certificate"

	default:
		http.Error(w, "Invalid source", http.StatusBadRequest)
		return
	}

	chainInfo := cert.AnalyzeCertificateChain(certs)

	htmlOutput, err := html.GenerateHTML(chainInfo, title)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating HTML: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlOutput))
}