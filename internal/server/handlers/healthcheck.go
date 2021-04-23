package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)


type HealthCheck struct {
	Status string `json:"status"`
	Version string `json:"version"`
	ReleaseID string `json:"releaseID"`
	Notes []string `json:"notes"`
	OutPut string `json:"output"`
}
//HealthCheckHandler ...
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheckResponse := HealthCheck{
		Status:    "pass",
		Version:   "1",
		ReleaseID: "0.0.1",
		Notes:     []string{""},
		OutPut:    "pass",
	}

	responseBytes, err := json.Marshal(&healthCheckResponse)
	if err != nil {
		log.Printf("Error while marshalling healthcheck error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBytes)
	if err != nil {
		log.Printf("Error while marshalling healthcheck error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
