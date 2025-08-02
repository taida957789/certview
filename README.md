# CertView 🔒

A comprehensive certificate analysis tool written in Go that can analyze X.509 certificates from files, certificate chains, or live domains via TLS/SNI handshake. The tool provides detailed certificate information, validates certificate chains, detects cross-signing relationships, and generates beautiful HTML reports with embedded CSS styling.

## Features

- **Multiple Input Sources**:
  - Certificate files (PEM/DER format)
  - Certificate chains
  - Live domain analysis via TLS/SNI handshake
  
- **Comprehensive Analysis**:
  - Certificate structure explanation
  - All certificate fields and extensions
  - Complete signing chain analysis
  - Cross-signing detection and visualization
  - Certificate validation and expiry checking
  
- **Rich Output**:
  - HTML reports with embedded CSS styling
  - Interactive certificate chain visualization
  - Detailed field breakdowns
  - Extension analysis with critical flag detection
  
- **Dual Operation Modes**:
  - CLI tool for command-line usage
  - Standalone HTTP server with web interface

## Installation

### From Source

```bash
git clone <repository-url>
cd certview
go build -o certview
```

### Direct Build

```bash
go build -o certview
```

## Usage

### CLI Mode

#### Analyze a certificate file:
```bash
./certview certificate.pem
./certview certificate.crt
./certview certificate.der
```

#### Analyze a domain (live TLS handshake):
```bash
./certview google.com:443
./certview example.com  # defaults to port 443
```

#### Output HTML to file:
```bash
./certview google.com:443 > analysis.html
```

### Server Mode

Start the web server:
```bash
./certview -server
./certview -server -port=8080
```

Then open your browser to `http://localhost:8080` to access the web interface.

The web interface supports:
- Domain analysis with live TLS handshake
- Certificate file upload (PEM/DER formats)
- Paste certificate data directly

## Examples

### CLI Examples

```bash
# Analyze Google's certificate
./certview google.com:443

# Analyze a local certificate file
./certview /path/to/certificate.pem

# Analyze with custom port
./certview example.com:8443

# Save analysis to HTML file
./certview google.com > google-cert-analysis.html
```

### Server Mode

```bash
# Start server on default port 8080
./certview -server

# Start server on custom port
./certview -server -port=9000
```

## Output Features

The generated HTML report includes:

- **Certificate Chain Visualization**: Visual representation of the certificate chain with CA certificates, end-entity certificates, expired certificates, and cross-signed certificates clearly marked
- **Validation Status**: Chain validation results with detailed error reporting
- **Cross-Signing Detection**: Identification and visualization of cross-signing relationships
- **Detailed Certificate Information**:
  - Subject and Issuer information
  - Serial numbers and validity periods
  - Public key algorithms and sizes
  - Signature algorithms
  - Key usage and extended key usage
  - Subject Alternative Names (SANs)
  - All certificate extensions with OIDs and critical flags
- **Interactive Interface**: Expandable certificate sections with toggle functionality
- **Responsive Design**: Mobile-friendly interface with modern CSS styling

## Architecture

```
certview/
├── main.go                 # Entry point
├── cmd/
│   ├── cli.go             # CLI command handling
│   └── server.go          # HTTP server implementation
├── pkg/
│   ├── cert/
│   │   ├── parser.go      # Certificate file parsing
│   │   ├── fetcher.go     # TLS handshake & cert retrieval
│   │   └── analyzer.go    # Certificate analysis & validation
│   └── html/
│       ├── generator.go   # HTML output generation
│       └── templates.go   # HTML templates with CSS
└── README.md
```

## Supported Formats

- **Certificate Files**: PEM (.pem, .crt, .cer), DER (.der)
- **Certificate Chains**: Multiple certificates in single PEM file
- **Live Domains**: Any domain with TLS enabled
- **Output**: HTML with embedded CSS (no external dependencies)

## Security Features

- **Certificate Validation**: Complete chain validation with detailed error reporting
- **Expiry Detection**: Clear indication of expired certificates
- **Cross-Signing Analysis**: Detection of multiple signing paths
- **Extension Analysis**: Detailed breakdown of all certificate extensions
- **Critical Flag Detection**: Identification of critical vs non-critical extensions

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **"no valid certificates found"**: Ensure the certificate file is in PEM or DER format
2. **"failed to connect"**: Check that the domain and port are correct and accessible
3. **"certificate validation failed"**: This is informational - the tool will still analyze invalid chains
4. **File upload issues**: Ensure certificate files are under 10MB and in supported formats

### Debugging

Use the CLI mode for detailed error messages:
```bash
./certview problematic-cert.pem
```

The server mode provides user-friendly error messages in the web interface.