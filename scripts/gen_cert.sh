#bin/sh

DOMAIN="$1"
SERIAL="$2"

cat > cert.conf <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage=digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName=@alt_names

[alt_names]
DNS.1 = $DOMAIN
EOF

cat > csr.conf <<EOF
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
CN = $DOMAIN

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $DOMAIN
EOF

openssl req -new -key cert.key -sha256 -config csr.conf | openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -set_serial "$2" -out /certs/"$1.crt" -extfile cert.conf 



