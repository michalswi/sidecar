package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Handlers struct {
	logger *log.Logger
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func main() {

	logger := log.New(os.Stdout, "web app ", log.LstdFlags|log.Lshortfile)
	s := NewHandlers(logger)
	logger.Println("Server is starting...")

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("PORT env var is missing.")
	}

	pxport := os.Getenv("PPORT")
	if pxport == "" {
		logger.Fatal("PPORT env var is missing.")
	}

	// addr := net.JoinHostPort("0.0.0.0", port)
	// fmt.Println(addr)

	// WEBAPP
	rserver := http.NewServeMux()

	rserver.Handle("/", s.index())
	// OR
	// if I use 'http.HandleFunc()' instead of 'rserver.HandleFunc()' it won't work as expected
	// rserver.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Fprintf(w, "Hi world!\n")
	// })

	go func() {
		logger.Println("Webapp is ready to handle requests at port", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", port), rserver); err != nil {
			logger.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	// PROXY
	// parse, without 'http://' - unsupported protocol scheme "localhost"
	url, err := url.Parse("http://localhost:" + port)
	if err != nil {
		logger.Fatalf("Error parsing backend url: %v", err)
	}

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	rproxy := http.NewServeMux()

	rproxy.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		logger.Printf("Proxy request dump: %v", req)
		proxy.ServeHTTP(rw, req)
	})

	logger.Println("Proxy is ready to handle requests at port", pxport)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", pxport), rproxy); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", pxport, err)
	}

	logger.Printf("Server stopped")
}

func (h *Handlers) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)

		h.logger.Printf("--> %s %s", r.Method, r.URL.Path)
		fmt.Fprintln(w, "Cloud runner!")
	})
}
