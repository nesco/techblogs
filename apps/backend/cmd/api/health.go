package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthAPI struct {
	startTime time.Time
}

func NewHealthAPI(startTime time.Time) *HealthAPI {
	return &HealthAPI{
		startTime,
	}
}

func (a *HealthAPI) Read(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().UTC().Format(time.RFC3339),
		"uptime_s":  time.Since(a.startTime).Seconds(),
	}
	json.NewEncoder(w).Encode(resp)
}
