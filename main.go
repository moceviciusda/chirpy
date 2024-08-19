package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("fileServerHits reset successfully"))
}

func main() {
	const port = "8080"
	config := apiConfig{}

	handler := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	handler.Handle("/app/", config.middlewareMetricsInc(fileServerHandler))

	handler.HandleFunc("GET /admin/metrics", config.metricsHandler)
	handler.HandleFunc("/admin/reset", config.resetHandler)

	handler.HandleFunc("GET /api/healthz", healthz)
	handler.HandleFunc("POST /api/validate_chirp", validateChirp)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: handler,
	}

	server.ListenAndServe()
}
