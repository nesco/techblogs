package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nesco/techblogs/backend/internal/blogs"
)

type BlogsAPI struct {
	repo *blogs.Repository
}

func blogsDataToCards(items []blogs.BlogInfo) (string, error) {
	var buffer bytes.Buffer
	if err := blogs.BlogListTemplate.Execute(&buffer, items); err != nil {
		return "", fmt.Errorf("error parsing blog entries: %w", err)
	}
	return buffer.String(), nil
}

func blogsDataToFeed(items []blogs.BlogInfo) (string, error) {
	var buffer bytes.Buffer
	if err := blogs.BlogFeedTemplate.Execute(&buffer, items); err != nil {
		return "", fmt.Errorf("error parsing blog entries: %w", err)
	}
	return buffer.String(), nil
}

func encodeBlogsJSON(w http.ResponseWriter, items []blogs.BlogInfo) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func encodeBlogsHTML(w http.ResponseWriter, items []blogs.BlogInfo) {
	content, err := blogsDataToCards(items)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, content)
}

func encodeBlogsRSS(w http.ResponseWriter, items []blogs.BlogInfo) {
	content, err := blogsDataToFeed(items)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	io.WriteString(w, content)
}

func NewBlogsAPI(repo *blogs.Repository) *BlogsAPI {
	return &BlogsAPI{repo: repo}
}

func (a *BlogsAPI) Read(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	var kind blogs.Kind
	var items []blogs.BlogInfo
	var err error

	w.Header().Set("Vary", "Accept")

	if collection != "" {
		var ok bool
		kind, ok = blogs.KindByCollection[collection]
		if !ok {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}
		items, err = a.repo.GetBlogsByKind(kind)
	} else {
		items, err = a.repo.GetAllBlogs()
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "application/json") {
		encodeBlogsJSON(w, items)
		return
	}

	if strings.Contains(accept, "application/rss+xml") {
		encodeBlogsRSS(w, items)
		return
	}

	// Default to HTML for text/html, */* or empty Accept header
	if accept == "" || accept == "*/*" || strings.Contains(accept, "text/html") {
		encodeBlogsHTML(w, items)
		return
	}

	http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
}

func (a *BlogsAPI) RSS(w http.ResponseWriter, r *http.Request) {
	items, err := a.repo.GetAllBlogs()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	encodeBlogsRSS(w, items)
}
