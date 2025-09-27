#!/bin/bash
set -e

SSL_DIR=./ssl
mkdir -p $SSL_DIR

# 1. สร้าง Root CA (ครั้งแรกครั้งเดียว)
if [ ! -f $SSL_DIR/LocalDevRootCA.key ]; then
  echo "Generating Local Root CA..."
  openssl genrsa -out $SSL_DIR/LocalDevRootCA.key 4096
  openssl req -x509 -new -nodes -sha256 -days 3650 \
    -key $SSL_DIR/LocalDevRootCA.key \
    -out $SSL_DIR/LocalDevRootCA.crt \
    -subj "/C=TH/ST=Bangkok/L=Bangkok/O=Local Dev/OU=Root CA/CN=LocalDevRootCA"
fi

# 2. สร้าง server key + CSR
openssl genrsa -out $SSL_DIR/hospital-a.api.co.th.key 2048
openssl req -new -sha256 \
  -key $SSL_DIR/hospital-a.api.co.th.key \
  -out $SSL_DIR/hospital-a.api.co.th.csr \
  -subj "/C=TH/ST=Bangkok/L=Bangkok/O=Hospital A/OU=Dev/CN=hospital-a.api.co.th"

# 3. เขียนไฟล์ SAN
cat > $SSL_DIR/hospital-a.api.co.th.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
subjectAltName=@alt_names

[alt_names]
DNS.1=hospital-a.api.co.th
EOF

# 4. ใช้ Root CA เซ็นออก cert
openssl x509 -req -sha256 -days 825 \
  -in $SSL_DIR/hospital-a.api.co.th.csr \
  -CA $SSL_DIR/LocalDevRootCA.crt -CAkey $SSL_DIR/LocalDevRootCA.key -CAcreateserial \
  -out $SSL_DIR/hospital-a.api.co.th.crt \
  -extfile $SSL_DIR/hospital-a.api.co.th.ext

echo "✅ Certificate generated!"
echo "  Server Cert : $SSL_DIR/hospital-a.api.co.th.crt"
echo "  Server Key  : $SSL_DIR/hospital-a.api.co.th.key"
echo "  Root CA     : $SSL_DIR/LocalDevRootCA.crt (import into Windows trust store)"
