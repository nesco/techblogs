package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthAPI struct {
	startTime time.Time
}

type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	UptimeS   float64 `json:"uptime_s"`
}

func NewHealthAPI(startTime time.Time) *HealthAPI {
	return &HealthAPI{
		startTime,
	}
}

func (a *HealthAPI) Read(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		UptimeS:   time.Since(a.startTime).Seconds(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
