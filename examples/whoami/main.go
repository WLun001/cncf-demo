package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	port string
	name string
)

func init() {
	flag.StringVar(&port, "port", getEnv("WHOAMI_PORT_NUMBER", "3000"), "give me a port number")
	flag.StringVar(&name, "name", os.Getenv("WHOAMI_NAME"), "give me a name")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", handle(whoamiHandler, true))
	log.Printf("Starting up on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))

}

func handle(next http.HandlerFunc, verbose bool) http.Handler {
	if !verbose {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		// <remote_IP_address> - [<timestamp>] "<request_method> <request_path> <request_protocol>" -
		log.Printf("%s - - [%s] \"%s %s %s\" - -", getIP(r), time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto)
	})
}

func whoamiHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	wait := u.Query().Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	if name != "" {
		_, _ = fmt.Fprintln(w, "Name:", name)
	}

	hostname, _ := os.Hostname()
	_, _ = fmt.Fprintln(w, "Hostname:", hostname)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			_, _ = fmt.Fprintln(w, "IP:", ip)
		}
	}

	_, _ = fmt.Fprintln(w, "RemoteAddr:", req.RemoteAddr)

	ipAddr := getIP(req)
	ipapiClient := http.Client{}
	ipapiReq, err := http.NewRequest("GET", fmt.Sprintf("https://ipapi.co/%s/json", ipAddr), nil)
	ipapiReq.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
	resp, err := ipapiClient.Do(ipapiReq)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintln(w, "ipapi:", string(body))

	if err := req.Write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getIP(req *http.Request) string {
	forwardedFor := req.Header.Get("X-Forwarded-For")
	realIP := req.Header.Get("X-Real-Ip")
	ipAddress := req.RemoteAddr
	if forwardedFor != "" {
		ipAddress = forwardedFor

		ips := strings.Split(forwardedFor, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	} else if realIP != "" {
		ipAddress = realIP
	}

	return ipAddress
}
