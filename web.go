package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// PORT=8080 go run web.go

func main() {
	logger := log.New(os.Stdout, "web app ", log.LstdFlags|log.Lshortfile)
	logger.Println("Server is starting...")

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("PORT env var is missing.")
	}

	rserver := http.NewServeMux()

	rserver.Handle("/", index())

	logger.Println("Web server is ready at", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), rserver); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", port, err)
	}

	logger.Printf("Server stopped")
}

func index() http.Handler {
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

		log.Printf("--> %s %s", r.Method, r.URL.Path)
		fmt.Fprintln(w, "Cloud runner!")
	})
}
