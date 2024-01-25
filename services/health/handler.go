package services

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthCheck struct {
	Status string
}

func Health(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(HealthCheck{Status: "OK"})

	if err != nil {
		log.Panicln(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
