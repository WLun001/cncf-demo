package main

import (
	"crypto/tls"
	"crypto/x509"
	"go-mtls/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling request from %s\n", r.RemoteAddr)
		_, err := w.Write([]byte("Hello from the TLS-enabled server!\n"))
		if err != nil {
			log.Println(err)
			return
		}
	})

	caCert, err := os.ReadFile(utils.EnvDefault("CA_CERTS_PATH", "certs/ca.crt"))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load server's certificate and key
	serverCert, err := tls.LoadX509KeyPair(
		utils.EnvDefault("CERT_PATH", "certs/server.crt"), utils.EnvDefault("KEY_PATH", "certs/server.key"))
	if err != nil {
		log.Fatal(err)
	}

	// Configure TLS with mTLS requirements
	config := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCert},
	}
	server := &http.Server{
		Addr:      utils.EnvDefault("SERVER_ADDR", ":8443"),
		TLSConfig: config,
	}

	log.Printf("Listening on %s\n", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
