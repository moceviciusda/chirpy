package main

import (
	"net/http"
)

func main() {
	const port = "8080"

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: http.NewServeMux(),
	}

	server.ListenAndServe()
}
