package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/moceviciusda/chirpy/internal/database"
)

type apiConfig struct {
	db             *database.DB
	JWT_SECRET     string
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	host := os.Getenv("HOST")
	dbPath := os.Getenv("DB_PATH")

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize DB")
	}

	config := apiConfig{db, jwtSecret, 0}

	handler := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	handler.Handle("/app/", config.middlewareMetricsInc(fileServerHandler))

	handler.HandleFunc("GET /admin/metrics", config.metricsHandler)
	handler.HandleFunc("/admin/reset", config.resetHandler)

	handler.HandleFunc("GET /api/healthz", healthz)

	handler.HandleFunc("POST /api/login", config.login)

	handler.HandleFunc("POST /api/users", config.postUser)

	handler.HandleFunc("POST /api/chirps", config.postChirp)
	handler.HandleFunc("GET /api/chirps", config.getChirps)
	handler.HandleFunc("GET /api/chirps/{chirpID}", config.getChirpById)

	server := &http.Server{
		Addr:    host,
		Handler: handler,
	}

	server.ListenAndServe()
}
