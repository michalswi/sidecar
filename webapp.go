package main

import (
	"fmt"
	"log"
	"net/http"
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
	logger := log.New(os.Stdout, "webapp ", log.LstdFlags|log.Lshortfile)
	s := NewHandlers(logger)
	logger.Println("Server is starting...")

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("PORT variable is missing.")
	}

	rserver := http.NewServeMux()

	rserver.Handle("/", s.index())

	logger.Println("Server is ready to handle requests at port", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), rserver); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", port, err)
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
