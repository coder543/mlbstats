package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Identifiable struct {
	ID           int    `json:"id"`
	Name         string `json:"name,omitempty"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation,omitempty"`
}

type ErrResponse struct {
	Reason string `json:"reason"`
}

func WriteErrResponse(w http.ResponseWriter, code int, reason string) {
	w.WriteHeader(code)
	marshalled, err := json.Marshal(ErrResponse{Reason: reason})
	if err != nil {
		log.Printf("Unable to write err response: %v", err)
		return
	}

	_, err = w.Write(marshalled)
	if err != nil {
		log.Printf("Unable to write err response: %v", err)
		return
	}
}
