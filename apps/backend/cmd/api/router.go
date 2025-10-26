package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/nesco/techblogs/backend/internal/blogs"
)

func registerRoutes(mux *http.ServeMux, startTime time.Time, db *sql.DB) {
	healthAPI := NewHealthAPI(startTime)
	blogsRepo := blogs.NewRepository(db)
	blogsAPI := NewBlogsAPI(blogsRepo)

	mux.HandleFunc("GET /health", healthAPI.Read)
	mux.HandleFunc("GET /blogs", blogsAPI.Read)
	mux.HandleFunc("GET /blogs/rss.xml", blogsAPI.RSS)
	mux.HandleFunc("GET /blogs/{collection}", blogsAPI.Read)
}
