package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	logger := log.New(os.Stdout, "proxy ", log.LstdFlags|log.Lshortfile)
	logger.Println("Proxy is starting...")

	// Proxy Port
	proxyPort := os.Getenv("PPORT")
	if proxyPort == "" {
		logger.Fatal("PPORT env var is missing.")
	}

	// Application Port
	appPort := os.Getenv("APORT")
	if appPort == "" {
		logger.Fatal("APORT env var is missing.")
	}

	// Application IP
	appIP := getEnv("AIP", "localhost")
	if appIP == "" {
		logger.Fatal("AIP env var is missing.")
	}

	url, err := url.Parse("http://" + appIP + ":" + appPort)
	if err != nil {
		logger.Fatalf("Error parsing backend url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	rproxy := http.NewServeMux()

	rproxy.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		logger.Printf("request dump: %v", req)
		proxy.ServeHTTP(rw, req)
	})

	logger.Println("Proxy is ready to handle requests at port", proxyPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", proxyPort), rproxy); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", proxyPort, err)
	}

	logger.Printf("Server stopped")
}
