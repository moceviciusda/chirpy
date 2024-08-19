package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type reqBody struct {
	Body string `json:"body"`
}

type errorBody struct {
	Error string `json:"error"`
}

type cleanedBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := reqBody{}

	w.Header().Set("Content-Type", "application/json")

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

	respBody := cleanedBody{strings.Join(words, " ")}
	data, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Write(data)
}
