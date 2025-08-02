#!/bin/bash

# Generate Cross-Signed Certificate Chain for example.com
# This script creates two root CAs and demonstrates cross-signing

set -e

echo "🔐 Generating Cross-Signed Certificate Chain for example.com"
echo "============================================================"

# Create directories
mkdir -p cross_sign_demo/{ca1,ca2,intermediate,end_entity}
cd cross_sign_demo

# Generate Root CA 1
echo "📜 Step 1: Creating Root CA 1 (TechCorp Root CA)"
echo "================================================"

# Root CA 1 private key
openssl genrsa -out ca1/root_ca1.key 4096

# Root CA 1 certificate
cat > ca1/root_ca1.conf << EOF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
prompt = no

[req_distinguished_name]
C = US
ST = California
L = San Francisco
O = TechCorp
OU = Certificate Authority
CN = TechCorp Root CA

[v3_ca]
basicConstraints = critical,CA:TRUE
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
EOF

openssl req -new -x509 -days 7300 -key ca1/root_ca1.key -out ca1/root_ca1.crt -config ca1/root_ca1.conf -extensions v3_ca

echo "✅ Root CA 1 created: TechCorp Root CA"

# Generate Root CA 2
echo "📜 Step 2: Creating Root CA 2 (GlobalTrust Root CA)"
echo "=================================================="

# Root CA 2 private key
openssl genrsa -out ca2/root_ca2.key 4096

# Root CA 2 certificate
cat > ca2/root_ca2.conf << EOF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
prompt = no

[req_distinguished_name]
C = GB
ST = London
L = London
O = GlobalTrust
OU = Root Certificate Authority
CN = GlobalTrust Root CA

[v3_ca]
basicConstraints = critical,CA:TRUE
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
EOF

openssl req -new -x509 -days 7300 -key ca2/root_ca2.key -out ca2/root_ca2.crt -config ca2/root_ca2.conf -extensions v3_ca

echo "✅ Root CA 2 created: GlobalTrust Root CA"

# Generate Intermediate CA key and CSR
echo "📜 Step 3: Creating Intermediate CA (SecureNet Intermediate CA)"
echo "============================================================="

# Intermediate CA private key
openssl genrsa -out intermediate/intermediate_ca.key 2048

# Intermediate CA CSR
cat > intermediate/intermediate_ca.conf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_intermediate_ca
prompt = no

[req_distinguished_name]
C = US
ST = New York
L = New York
O = SecureNet
OU = Intermediate Certificate Authority
CN = SecureNet Intermediate CA

[v3_intermediate_ca]
basicConstraints = critical,CA:TRUE,pathlen:0
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
EOF

openssl req -new -key intermediate/intermediate_ca.key -out intermediate/intermediate_ca.csr -config intermediate/intermediate_ca.conf -reqexts v3_intermediate_ca

echo "✅ Intermediate CA CSR created"

# Sign Intermediate CA with Root CA 1
echo "📜 Step 4: Cross-Signing - Intermediate signed by Root CA 1"
echo "=========================================================="

cat > intermediate/intermediate_ext1.conf << EOF
[v3_intermediate_ca]
basicConstraints = critical,CA:TRUE,pathlen:0
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
EOF

openssl x509 -req -in intermediate/intermediate_ca.csr -CA ca1/root_ca1.crt -CAkey ca1/root_ca1.key -CAcreateserial -out intermediate/intermediate_ca_signed_by_ca1.crt -days 3650 -extensions v3_intermediate_ca -extfile intermediate/intermediate_ext1.conf

echo "✅ Intermediate CA signed by Root CA 1"

# Sign Intermediate CA with Root CA 2 (CROSS-SIGNING!)
echo "📜 Step 5: Cross-Signing - Same Intermediate signed by Root CA 2"
echo "================================================================"

cat > intermediate/intermediate_ext2.conf << EOF
[v3_intermediate_ca]
basicConstraints = critical,CA:TRUE,pathlen:0
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
EOF

openssl x509 -req -in intermediate/intermediate_ca.csr -CA ca2/root_ca2.crt -CAkey ca2/root_ca2.key -CAcreateserial -out intermediate/intermediate_ca_signed_by_ca2.crt -days 3650 -extensions v3_intermediate_ca -extfile intermediate/intermediate_ext2.conf

echo "🔗 CROSS-SIGNING ACHIEVED! Same intermediate CA now signed by two different Root CAs"

# Generate End Entity Certificate for example.com
echo "📜 Step 6: Creating End Entity Certificate for example.com"
echo "========================================================="

# End entity private key
openssl genrsa -out end_entity/example.com.key 2048

# End entity CSR
cat > end_entity/example.com.conf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_end_entity
prompt = no

[req_distinguished_name]
C = US
ST = California
L = San Francisco
O = Example Corporation
OU = IT Department
CN = example.com

[v3_end_entity]
basicConstraints = critical,CA:FALSE
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth,clientAuth
subjectKeyIdentifier = hash
subjectAltName = @alt_names

[alt_names]
DNS.1 = example.com
DNS.2 = www.example.com
DNS.3 = api.example.com
EOF

openssl req -new -key end_entity/example.com.key -out end_entity/example.com.csr -config end_entity/example.com.conf -reqexts v3_end_entity

# Sign end entity with intermediate CA (using the one signed by CA1)
cat > end_entity/end_entity_ext.conf << EOF
[v3_end_entity]
basicConstraints = critical,CA:FALSE
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth,clientAuth
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
subjectAltName = @alt_names

[alt_names]
DNS.1 = example.com
DNS.2 = www.example.com
DNS.3 = api.example.com
EOF

openssl x509 -req -in end_entity/example.com.csr -CA intermediate/intermediate_ca_signed_by_ca1.crt -CAkey intermediate/intermediate_ca.key -CAcreateserial -out end_entity/example.com.crt -days 365 -extensions v3_end_entity -extfile end_entity/end_entity_ext.conf

echo "✅ End entity certificate created for example.com"

# Create certificate chains
echo "📜 Step 7: Creating Certificate Chains"
echo "====================================="

# Chain 1: example.com -> intermediate (signed by CA1) -> CA1
cat end_entity/example.com.crt intermediate/intermediate_ca_signed_by_ca1.crt ca1/root_ca1.crt > chain1_ca1_path.pem

# Chain 2: example.com -> intermediate (signed by CA2) -> CA2
cat end_entity/example.com.crt intermediate/intermediate_ca_signed_by_ca2.crt ca2/root_ca2.crt > chain2_ca2_path.pem

# Cross-signed demonstration chain: Contains both cross-signed intermediates
cat end_entity/example.com.crt intermediate/intermediate_ca_signed_by_ca1.crt intermediate/intermediate_ca_signed_by_ca2.crt ca1/root_ca1.crt ca2/root_ca2.crt > cross_signed_demo_chain.pem

echo "✅ Certificate chains created:"
echo "   - chain1_ca1_path.pem (example.com -> Intermediate[CA1] -> Root CA1)"
echo "   - chain2_ca2_path.pem (example.com -> Intermediate[CA2] -> Root CA2)"  
echo "   - cross_signed_demo_chain.pem (Contains cross-signed certificates)"

echo ""
echo "🎉 CROSS-SIGNED CERTIFICATE CHAIN GENERATION COMPLETE!"
echo "======================================================"
echo ""
echo "📊 Summary:"
echo "- Created 2 Root CAs: TechCorp Root CA & GlobalTrust Root CA"
echo "- Created 1 Intermediate CA: SecureNet Intermediate CA"
echo "- CROSS-SIGNED the Intermediate CA with both Root CAs"
echo "- Created end entity certificate for example.com"
echo "- Generated demonstration chains showing cross-signing"
echo ""
echo "🔍 Files created:"
echo "- cross_signed_demo_chain.pem (Best file to test cross-signing detection)"
echo "- chain1_ca1_path.pem (Traditional chain via CA1)"
echo "- chain2_ca2_path.pem (Traditional chain via CA2)"
echo ""
echo "✨ Ready to test with: ./certview cross_sign_demo/cross_signed_demo_chain.pem"

cd ..
echo ""
echo "🚀 Testing cross-signing detection with generated certificates..."
echo "================================================================"

# Test with our cross-signed chain
./certview cross_sign_demo/cross_signed_demo_chain.pem > cross_signed_test_report.html

echo "✅ Cross-signing test report generated: cross_signed_test_report.html"
echo ""
echo "📋 To view the results:"
echo "   open cross_signed_test_report.html"
echo ""
echo "🔍 The report should show cross-signing detection for the intermediate CA!"