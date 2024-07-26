
SERVER_CN=localhost

# CA
# 1. Generate CA's private key
openssl genrsa -out ca.key 2048

# 2. Create a self-signed CA certificate
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/O=${SERVER_CN}"


#####
# Server
# 1. Generate server's private key
openssl genrsa -out server.key 2048

# 2. Create a certificate signing request (CSR) for the server
openssl req -new -key server.key -out server.csr -subj "/O=${SERVER_CN}" -subj "/CN=${SERVER_CN}" -addext "subjectAltName = DNS:${SERVER_CN}"

# 3. Sign the CSR with the CA's key to create the server certificate
openssl x509 -req -extfile <(printf "subjectAltName=DNS:${SERVER_CN}") -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365


### client
# 1. Generate client's private key
openssl genrsa -out client.key 2048

# 2. Create a CSR for the client
openssl req -new -key client.key -out client.csr -subj "/O=${SERVER_CN}" -subj "/CN=${SERVER_CN}" -addext "subjectAltName = DNS:${SERVER_CN}"

# 3. Sign the CSR with the CA's key to create the client certificate
openssl x509 -req -extfile <(printf "subjectAltName=DNS:${SERVER_CN}") -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365

