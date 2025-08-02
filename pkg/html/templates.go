package html

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Certificate Analysis - {{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }

        .header {
            background: rgba(255, 255, 255, 0.95);
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
            text-align: center;
        }

        .header h1 {
            color: #4a5568;
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .header .subtitle {
            color: #718096;
            font-size: 1.2em;
        }

        .chain-overview {
            background: rgba(255, 255, 255, 0.95);
            padding: 25px;
            border-radius: 15px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }

        .chain-status {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
        }

        .status-badge {
            padding: 8px 16px;
            border-radius: 20px;
            font-weight: bold;
            margin-right: 15px;
        }

        .status-valid {
            background: #48bb78;
            color: white;
        }

        .status-invalid {
            background: #f56565;
            color: white;
        }

        .chain-visualization {
            margin: 20px 0;
        }

        .cert-chain {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }

        .cert-link {
            display: flex;
            align-items: center;
            justify-content: center;
            position: relative;
        }

        .cert-box {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 15px 25px;
            border-radius: 10px;
            min-width: 300px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
            position: relative;
        }

        .cert-box.ca {
            background: linear-gradient(135deg, #48bb78, #38a169);
        }

        .cert-box.expired {
            background: linear-gradient(135deg, #f56565, #e53e3e);
        }

        .cert-box.cross-signed {
            border: 3px dashed #fbb6ce;
        }

        .cert-arrow {
            color: #4a5568;
            font-size: 2em;
            margin: 5px 0;
            text-align: center;
            display: flex;
            justify-content: center;
            align-items: center;
        }

        .cert-details {
            background: rgba(255, 255, 255, 0.95);
            margin: 30px 0;
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }

        .cert-header {
            background: linear-gradient(135deg, #4a5568, #2d3748);
            color: white;
            padding: 20px;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .cert-header:hover {
            background: linear-gradient(135deg, #2d3748, #1a202c);
        }

        .cert-body {
            padding: 25px;
            display: none;
        }

        .cert-body.expanded {
            display: block;
        }

        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 25px;
        }

        .info-card {
            background: #f7fafc;
            padding: 20px;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }

        .info-card h4 {
            color: #4a5568;
            margin-bottom: 10px;
            font-size: 1.1em;
        }

        .info-card p {
            color: #718096;
            word-break: break-all;
        }

        .extensions-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        .extensions-table th,
        .extensions-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e2e8f0;
        }

        .extensions-table th {
            background: #edf2f7;
            font-weight: 600;
            color: #4a5568;
        }

        .extensions-table tr:hover {
            background: #f7fafc;
        }

        .critical {
            background: #fed7d7;
            color: #c53030;
            padding: 2px 6px;
            border-radius: 4px;
            font-size: 0.8em;
        }

        .errors-section {
            background: rgba(254, 215, 215, 0.9);
            border: 1px solid #fc8181;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
        }

        .errors-section h3 {
            color: #c53030;
            margin-bottom: 15px;
        }

        .error-list {
            list-style: none;
        }

        .error-list li {
            color: #c53030;
            margin-bottom: 5px;
            padding-left: 20px;
            position: relative;
        }

        .error-list li:before {
            content: "⚠";
            position: absolute;
            left: 0;
        }

        .cross-signing-section {
            background: rgba(255, 235, 230, 0.9);
            border: 1px solid #fbb6ce;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
        }

        .cross-signing-section h3 {
            color: #97266d;
            margin-bottom: 15px;
        }

        .toggle-icon {
            transition: transform 0.3s ease;
        }

        .toggle-icon.rotated {
            transform: rotate(180deg);
        }

        @media (max-width: 768px) {
            .container {
                padding: 10px;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .info-grid {
                grid-template-columns: 1fr;
            }
            
            .cert-box {
                min-width: 250px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔒 Certificate Analysis</h1>
            <div class="subtitle">{{.Title}}</div>
        </div>

        <div class="chain-overview">
            <div class="chain-status">
                {{if .ChainInfo.IsValid}}
                <div class="status-badge status-valid">✓ Valid Chain</div>
                {{else}}
                <div class="status-badge status-invalid">✗ Invalid Chain</div>
                {{end}}
                <span>{{len .ChainInfo.Certificates}} certificate(s) in chain</span>
            </div>

            {{if not .ChainInfo.IsValid}}
            <div class="errors-section">
                <h3>Chain Validation Errors</h3>
                <ul class="error-list">
                    {{range .ChainInfo.Errors}}
                    <li>{{.}}</li>
                    {{end}}
                </ul>
            </div>
            {{end}}

            {{if .ChainInfo.CrossSigning}}
            <div class="cross-signing-section">
                <h3>🔗 Cross-Signing Detected</h3>
                <p>Cross-signed certificates are identical certificates (same public key and subject) that have been signed by different Certificate Authorities, providing multiple validation paths.</p>
                {{range $key, $certs := .ChainInfo.CrossSigning}}
                <div style="margin: 15px 0; padding: 15px; background: rgba(255, 255, 255, 0.7); border-radius: 8px;">
                    <div style="font-weight: bold; color: #97266d; margin-bottom: 10px;">
                        📜 Certificate Group: {{len $certs}} cross-signed version(s)
                    </div>
                    <div style="font-size: 0.9em; color: #666; margin-bottom: 10px;">
                        {{$key}}
                    </div>
                    <div style="margin-left: 20px;">
                        {{range $i, $cert := $certs}}
                        <div style="margin: 5px 0; padding: 8px; background: rgba(151, 38, 109, 0.1); border-radius: 4px;">
                            <strong>Version {{add $i 1}}:</strong><br>
                            <span style="font-size: 0.85em;">
                                <strong>Subject:</strong> {{$cert.Subject}}<br>
                                <strong>Issuer:</strong> {{$cert.Issuer}}<br>
                                <strong>Serial:</strong> {{$cert.SerialNumber}}<br>
                                <strong>Valid:</strong> {{$cert.NotBefore.Format "2006-01-02"}} to {{$cert.NotAfter.Format "2006-01-02"}}
                            </span>
                        </div>
                        {{end}}
                    </div>
                </div>
                {{end}}
                <div style="margin-top: 15px; padding: 10px; background: rgba(72, 187, 120, 0.1); border-radius: 6px; font-size: 0.9em;">
                    💡 <strong>Why Cross-Signing Matters:</strong> Cross-signing provides redundancy and helps with certificate chain validation when different root stores are used across platforms and browsers.
                </div>
            </div>
            {{end}}

            <div class="chain-visualization">
                <h3>Certificate Chain Visualization</h3>
                {{if .ChainInfo.CrossSigning}}
                <div style="background: rgba(255, 248, 220, 0.9); padding: 20px; border-radius: 10px; margin-bottom: 20px;">
                    <h4 style="color: #b7791f; margin-bottom: 15px;">🔗 Cross-Signing Structure Detected</h4>
                    <p style="margin-bottom: 15px; color: #8b5e3c;">The visualization below shows the complex cross-signing relationships in your certificate chain:</p>
                    
                    <div class="cross-sign-tree">
                        {{range $key, $certs := .ChainInfo.CrossSigning}}
                        <div class="cross-sign-group" style="margin: 20px 0; padding: 15px; border: 2px dashed #d69e2e; border-radius: 10px; background: rgba(237, 242, 247, 0.5);">
                            <div style="text-align: center; margin-bottom: 15px;">
                                <div class="cert-box cross-signed ca" style="display: inline-block; margin: 0;">
                                    <strong>🔗 Cross-Signed Certificate</strong><br>
                                    {{(index $certs 0).Subject}}
                                </div>
                            </div>
                            
                            <div style="display: flex; justify-content: space-around; align-items: flex-start; flex-wrap: wrap; margin-top: 15px;">
                                {{range $i, $cert := $certs}}
                                <div style="margin: 10px; text-align: center; flex: 1; min-width: 250px;">
                                    <div class="cert-arrow" style="margin: 5px 0;">↑</div>
                                    <div class="cert-box ca" style="margin: 0;">
                                        <strong>🏛️ Issuer {{add $i 1}}:</strong><br>
                                        {{$cert.Issuer}}
                                    </div>
                                    <div style="font-size: 0.8em; color: #666; margin-top: 5px;">
                                        Serial: {{$cert.SerialNumber}}<br>
                                        Valid: {{$cert.NotBefore.Format "2006-01-02"}} to {{$cert.NotAfter.Format "2006-01-02"}}
                                    </div>
                                </div>
                                {{end}}
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
                
                {{if .ChainInfo.ChainPaths}}
                <div style="background: rgba(240, 253, 244, 0.9); padding: 20px; border-radius: 10px; margin: 20px 0;">
                    <h4 style="color: #2f855a; margin-bottom: 15px;">🛤️ Multiple Validation Paths</h4>
                    <p style="margin-bottom: 15px; color: #2d3748;">Due to cross-signing, this certificate can be validated through multiple paths:</p>
                    
                    {{range $pathIndex, $path := .ChainInfo.ChainPaths}}
                    <div style="margin: 15px 0; padding: 15px; background: rgba(255, 255, 255, 0.8); border-radius: 8px; border-left: 4px solid #48bb78;">
                        <h5 style="color: #2f855a; margin-bottom: 10px;">{{$path.Description}}</h5>
                        <div class="cert-chain" style="display: flex; flex-direction: column; align-items: center;">
                            {{range $i, $cert := $path.Path}}
                            <div class="cert-link">
                                <div class="cert-box{{if $cert.IsCA}} ca{{end}}{{if $cert.IsExpired}} expired{{end}}" style="max-width: 400px;">
                                    <strong>{{if $cert.IsCA}}🏛️ CA: {{else}}🌐 End Entity: {{end}}</strong><br>
                                    {{$cert.Subject}}
                                    {{if $cert.IsExpired}}<br><small>⚠️ EXPIRED</small>{{end}}
                                </div>
                            </div>
                            {{if ne $i (sub (len $path.Path) 1)}}
                            <div class="cert-arrow">↓</div>
                            {{end}}
                            {{end}}
                        </div>
                        <div style="margin-top: 10px; font-size: 0.85em; color: #4a5568;">
                            {{if $path.IsComplete}}
                            ✅ <strong>Complete validation path</strong>
                            {{else}}
                            ⚠️ <strong>Incomplete path</strong> - some intermediate certificates may be missing
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
                {{end}}
                {{end}}
                
                <div class="traditional-chain">
                    <h4 style="margin-bottom: 15px;">📋 All Certificates in Order</h4>
                    <div class="cert-chain">
                        {{range $i, $cert := .ChainInfo.Certificates}}
                        <div class="cert-link">
                            <div class="cert-box{{if $cert.IsCA}} ca{{end}}{{if $cert.IsExpired}} expired{{end}}{{if index $.ChainInfo.CrossSigning $cert.Subject}} cross-signed{{end}}">
                                <strong>{{if $cert.IsCA}}🏛️ CA: {{else}}🌐 End Entity: {{end}}</strong><br>
                                {{$cert.Subject}}
                                {{if $cert.IsExpired}}<br><small>⚠️ EXPIRED</small>{{end}}
                            </div>
                        </div>
                        {{if ne $i (sub (len $.ChainInfo.Certificates) 1)}}
                        <div class="cert-arrow">↓</div>
                        {{end}}
                        {{end}}
                    </div>
                </div>
            </div>
        </div>

        {{range $i, $cert := .ChainInfo.Certificates}}
        <div class="cert-details">
            <div class="cert-header" onclick="toggleCert({{$i}})">
                <h2>Certificate {{add $i 1}}: {{if $cert.IsCA}}Certificate Authority{{else}}End Entity{{end}}</h2>
                <span class="toggle-icon" id="icon-{{$i}}">▼</span>
            </div>
            <div class="cert-body" id="cert-{{$i}}">
                <div class="info-grid">
                    <div class="info-card">
                        <h4>Subject</h4>
                        <p>{{$cert.Subject}}</p>
                    </div>
                    <div class="info-card">
                        <h4>Issuer</h4>
                        <p>{{$cert.Issuer}}</p>
                    </div>
                    <div class="info-card">
                        <h4>Serial Number</h4>
                        <p>{{$cert.SerialNumber}}</p>
                    </div>
                    <div class="info-card">
                        <h4>Validity Period</h4>
                        <p><strong>Not Before:</strong> {{$cert.NotBefore.Format "2006-01-02 15:04:05 UTC"}}<br>
                        <strong>Not After:</strong> {{$cert.NotAfter.Format "2006-01-02 15:04:05 UTC"}}<br>
                        {{if $cert.IsExpired}}<span style="color: #c53030;">⚠️ EXPIRED</span>{{else}}<span style="color: #48bb78;">✓ Valid</span>{{end}}</p>
                    </div>
                    <div class="info-card">
                        <h4>Public Key</h4>
                        <p><strong>Algorithm:</strong> {{$cert.PublicKeyAlg}}<br>
                        {{if $cert.PublicKeySize}}<strong>Size:</strong> {{$cert.PublicKeySize}} bits{{end}}</p>
                    </div>
                    <div class="info-card">
                        <h4>Signature Algorithm</h4>
                        <p>{{$cert.SignatureAlg}}</p>
                    </div>
                    {{if $cert.KeyUsage}}
                    <div class="info-card">
                        <h4>Key Usage</h4>
                        <p>{{range $cert.KeyUsage}}{{.}}<br>{{end}}</p>
                    </div>
                    {{end}}
                    {{if $cert.ExtKeyUsage}}
                    <div class="info-card">
                        <h4>Extended Key Usage</h4>
                        <p>{{range $cert.ExtKeyUsage}}{{.}}<br>{{end}}</p>
                    </div>
                    {{end}}
                    {{if $cert.SANs}}
                    <div class="info-card">
                        <h4>Subject Alternative Names</h4>
                        <p>{{range $cert.SANs}}{{.}}<br>{{end}}</p>
                    </div>
                    {{end}}
                </div>

                {{if $cert.Extensions}}
                <h3>Extensions</h3>
                <table class="extensions-table">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>OID</th>
                            <th>Critical</th>
                            <th>Value</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range $cert.Extensions}}
                        <tr>
                            <td>{{.Name}}</td>
                            <td>{{.OID}}</td>
                            <td>{{if .Critical}}<span class="critical">CRITICAL</span>{{else}}No{{end}}</td>
                            <td style="word-break: break-all; max-width: 200px;">{{.Value}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{end}}
            </div>
        </div>
        {{end}}
    </div>

    <script>
        function toggleCert(index) {
            const body = document.getElementById('cert-' + index);
            const icon = document.getElementById('icon-' + index);
            
            if (body.classList.contains('expanded')) {
                body.classList.remove('expanded');
                icon.classList.remove('rotated');
            } else {
                body.classList.add('expanded');
                icon.classList.add('rotated');
            }
        }

        // Expand first certificate by default
        document.addEventListener('DOMContentLoaded', function() {
            toggleCert(0);
        });
    </script>
</body>
</html>`

const webFormTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CertView - Certificate Analysis Tool</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            padding: 40px;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 90%;
        }

        h1 {
            text-align: center;
            color: #4a5568;
            margin-bottom: 10px;
            font-size: 2.5em;
        }

        .subtitle {
            text-align: center;
            color: #718096;
            margin-bottom: 40px;
            font-size: 1.2em;
        }

        .form-group {
            margin-bottom: 30px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #4a5568;
        }

        input[type="text"], input[type="file"], textarea {
            width: 100%;
            padding: 12px;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s ease;
        }

        input[type="text"]:focus, input[type="file"]:focus, textarea:focus {
            outline: none;
            border-color: #667eea;
        }

        textarea {
            resize: vertical;
            min-height: 120px;
            font-family: monospace;
        }

        .radio-group {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
        }

        .radio-option {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        input[type="radio"] {
            width: 18px;
            height: 18px;
        }

        .input-section {
            display: none;
        }

        .input-section.active {
            display: block;
        }

        button {
            width: 100%;
            padding: 15px;
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 18px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s ease;
        }

        button:hover {
            transform: translateY(-2px);
        }

        button:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }

        .example {
            background: #f7fafc;
            padding: 10px;
            border-radius: 6px;
            font-family: monospace;
            font-size: 14px;
            color: #718096;
            margin-top: 5px;
        }

        .loading {
            display: none;
            text-align: center;
            margin-top: 20px;
        }

        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #667eea;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔒 CertView</h1>
        <div class="subtitle">Certificate Analysis Tool</div>

        <form id="certForm" method="POST" enctype="multipart/form-data">
            <div class="form-group">
                <label>Analysis Source:</label>
                <div class="radio-group">
                    <div class="radio-option">
                        <input type="radio" id="domain" name="source" value="domain" checked>
                        <label for="domain">Domain</label>
                    </div>
                    <div class="radio-option">
                        <input type="radio" id="file" name="source" value="file">
                        <label for="file">Certificate File</label>
                    </div>
                    <div class="radio-option">
                        <input type="radio" id="paste" name="source" value="paste">
                        <label for="paste">Paste Certificate</label>
                    </div>
                </div>
            </div>

            <div id="domain-section" class="input-section active">
                <div class="form-group">
                    <label for="domain-input">Domain and Port:</label>
                    <input type="text" id="domain-input" name="domain" placeholder="example.com:443">
                    <div class="example">Example: google.com:443 or just google.com (defaults to port 443)</div>
                </div>
            </div>

            <div id="file-section" class="input-section">
                <div class="form-group">
                    <label for="file-input">Certificate File:</label>
                    <input type="file" id="file-input" name="certfile" accept=".pem,.crt,.cer,.der">
                    <div class="example">Supports PEM and DER formats (.pem, .crt, .cer, .der)</div>
                </div>
            </div>

            <div id="paste-section" class="input-section">
                <div class="form-group">
                    <label for="paste-input">Paste Certificate:</label>
                    <textarea id="paste-input" name="certdata" placeholder="-----BEGIN CERTIFICATE-----
MIIFXzCCA0egAwIBAgIRAOJyQ...
-----END CERTIFICATE-----"></textarea>
                    <div class="example">Paste PEM formatted certificate(s) here</div>
                </div>
            </div>

            <button type="submit">🔍 Analyze Certificate</button>
        </form>

        <div class="loading" id="loading">
            <div class="spinner"></div>
            <p>Analyzing certificate...</p>
        </div>
    </div>

    <script>
        const sourceRadios = document.querySelectorAll('input[name="source"]');
        const sections = document.querySelectorAll('.input-section');
        
        sourceRadios.forEach(radio => {
            radio.addEventListener('change', function() {
                sections.forEach(section => section.classList.remove('active'));
                document.getElementById(this.value + '-section').classList.add('active');
            });
        });

        document.getElementById('certForm').addEventListener('submit', function(e) {
            const loading = document.getElementById('loading');
            const button = document.querySelector('button[type="submit"]');
            
            loading.style.display = 'block';
            button.disabled = true;
            button.textContent = 'Analyzing...';
        });
    </script>
</body>
</html>`