package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type postUserReq struct {
	Email string `json:"email"`
}

type postChirpReq struct {
	Body string `json:"body"`
}

type errorBody struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	reqBody := postUserReq{}

	err := decoder.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(500)
		respBody := errorBody{fmt.Sprintf("error decoding request body: %s", err)}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			return
		}

		w.Write(data)
		return
	}

	user, err := cfg.db.CreateUser(reqBody.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to save user: %s", err)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	reqBody := postChirpReq{}

	err := decoder.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(500)
		respBody := errorBody{fmt.Sprintf("error decoding request body: %s", err)}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			return
		}

		w.Write(data)
		return
	}

	if len(reqBody.Body) > 140 {
		w.WriteHeader(400)
		respBody := errorBody{"chirp is too long"}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}

		w.Write(data)
		return
	}

	words := strings.Fields(reqBody.Body)

	for i, w := range words {
		switch strings.ToLower(w) {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		}
	}

	body := strings.Join(words, " ")

	chirp, err := cfg.db.CreateChirp(body)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to save chirp: %s", err)
		return
	}

	data, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error getting chirps: %s", err)
		return
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Invalid chirp ID: %s", r.PathValue("chirpID"))
		return
	}

	chirp, err := cfg.db.GetChirp(chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	data, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
