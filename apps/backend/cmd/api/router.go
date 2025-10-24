package main

import (
	"net/http"
	"time"
)

func registerRoutes(mux *http.ServeMux, startTime time.Time) {
	healthAPI := NewHealthAPI(startTime)
	blogsAPI := NewBlogsAPI()

	mux.HandleFunc("GET /health", healthAPI.Read)
	mux.HandleFunc("GET /blogs", blogsAPI.Read)
	mux.HandleFunc("GET /blogs/{collection}", blogsAPI.Read)
}
