package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

	// Configure TLS with mTLS requirements
	config := &tls.Config{
		ClientCAs:      caCertPool,
		ClientAuth:     tls.RequireAndVerifyClientCert,
		GetCertificate: GetCertificate,
	}
	server := &http.Server{
		Addr:      utils.EnvDefault("SERVER_ADDR", ":8443"),
		TLSConfig: config,
	}

	log.Printf("Listening on %s\n", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(
		utils.EnvDefault("CERT_PATH", "certs/server.crt"),
		utils.EnvDefault("KEY_PATH", "certs/server.key"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not load TLS cert: %s", err)
	}
	parsedCert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("could not parsed TLS cert: %s", err)
	}
	log.Printf("Loaded certificate (%s): NotBefore: %s, NotAfter: %s", parsedCert.Version, parsedCert.NotBefore, parsedCert.NotAfter)
	return &cert, nil
}
