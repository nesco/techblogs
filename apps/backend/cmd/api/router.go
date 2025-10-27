package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/nesco/techblogs/backend/internal/blogs"
)

func registerRoutes(mux *http.ServeMux, startTime time.Time, db *sql.DB) {
	healthHandler := NewHealthHandler(startTime)
	blogsRepo := blogs.NewRepository(db)
	blogsHandler := NewBlogsHandler(blogsRepo)

	mux.HandleFunc("GET /health", healthHandler.Read)
	mux.HandleFunc("GET /blogs", blogsHandler.Read)
	mux.HandleFunc("GET /blogs/rss.xml", blogsHandler.RSS)
	mux.HandleFunc("GET /blogs/{collection}", blogsHandler.Read)
}
