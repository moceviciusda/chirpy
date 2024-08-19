package main

import (
	"net/http"
)

func main() {
	const port = "8080"

	handler := http.NewServeMux()

	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	handler.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: handler,
	}

	server.ListenAndServe()
}
