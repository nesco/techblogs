package main

import (
	"net/http"
	"time"
)

func registerRoutes(mux *http.ServeMux, startTime time.Time) {
	healtAPI := NewHealthAPI(startTime)
	blogsAPI := NewBlogsAPI()

	mux.HandleFunc("GET /health", healtAPI.Read)
	mux.HandleFunc("GET /blogs", blogsAPI.Read)
	mux.HandleFunc("GET /blogs/{kind}", blogsAPI.Read)
}
