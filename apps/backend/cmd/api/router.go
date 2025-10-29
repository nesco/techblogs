package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/nesco/techblogs/backend/internal/blogs"
	"github.com/nesco/techblogs/backend/internal/home"
	"go.uber.org/zap"
)

func registerRoutes(mux *http.ServeMux, startTime time.Time, db *sql.DB, logger *zap.SugaredLogger) {
	healthHandler := NewHealthHandler(startTime)
	blogsRepo := blogs.NewRepository(db)
	blogsHandler := NewBlogsHandler(blogsRepo)
	homeHandler := &home.HomeHandler{Logger: *logger, Repo: blogsRepo}

	// Home page
	mux.HandleFunc("GET /", homeHandler.Read)

	// API endpoints
	mux.HandleFunc("GET /api/health", healthHandler.Read)
	mux.HandleFunc("GET /api/blogs", blogsHandler.Read)
	mux.HandleFunc("GET /api/blogs/rss.xml", blogsHandler.RSS)
	mux.HandleFunc("GET /api/blogs/{collection}", blogsHandler.Read)
}
