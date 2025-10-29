// Package home renders the home page of techblo.gs
package home

import (
	"bytes"
	"io"
	"net/http"

	"github.com/nesco/techblogs/backend/internal/blogs"
	"go.uber.org/zap"
)

type HomeHandler struct {
	Logger zap.SugaredLogger
	Repo   *blogs.Repository
}

type HomePageData struct {
	People        []blogs.BlogInfo
	Organizations []blogs.BlogInfo
}

func (a *HomeHandler) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	people, err := a.Repo.GetBlogsByKind(blogs.Individual)
	if err != nil {
		a.Logger.Errorf("error fetching people blogs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	organizations, err := a.Repo.GetBlogsByKind(blogs.Organization)
	if err != nil {
		a.Logger.Errorf("error fetching organization blogs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := HomePageData{
		People:        people,
		Organizations: organizations,
	}

	var buffer bytes.Buffer
	if err := homePageTemplate.Execute(&buffer, data); err != nil {
		a.Logger.Errorf("error parsing main page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = io.WriteString(w, buffer.String())
	if err != nil {
		a.Logger.Errorw("error encoding home page response", "error", err, "client", r.Header.Get("X-Real-IP"))
	}
}
