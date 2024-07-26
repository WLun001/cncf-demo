package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go-mtls/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load client's certificate and key
	cert, err := tls.LoadX509KeyPair(
		utils.EnvDefault("CERT_PATH", "certs/client.crt"), utils.EnvDefault("KEY_PATH", "certs/client.key"))
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile(utils.EnvDefault("CA_CERTS_PATH", "certs/ca.crt"))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS with the client certificate
	config := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: transport}

	if utils.EnvDefault("LOOP", "false") == "true" {
		for {
			sendRequest(client)
			time.Sleep(60 * time.Second)
		}
	} else {
		sendRequest(client)
	}
}

func sendRequest(client *http.Client) {
	resp, err := client.Get(utils.EnvDefault("SERVER_URL", "https://localhost:8443"))
	if err != nil {
		log.Println(err)
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)

		log.Printf("Sending request to %s\n", resp.Request.Host)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Response: %s", body)
	}
}
