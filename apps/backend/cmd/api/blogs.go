package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/nesco/techblogs/backend/internal/blogs"
)

var blogEntriesTemplate = template.Must(template.New("BlogEntries").Parse(`
	{{- range . -}}
	<article class="card">
	 <h3><a href="{{ .BlogHref }}">{{ .BlogName }}</a></h3>
		{{- if .LatestArticleHref }}
		<p>Latest: <a href="{{ .LatestArticleHref }}">{{ if .LatestArticleName }}{{ .LatestArticleName }}{{ else }}{{ .LatestArticleHref }}{{ end }}</a></p>
		{{- end }}
	</article>
	{{- end -}}
`))

type BlogsAPI struct {
	repo *blogs.Repository
}

func blogsDataToCards(blogsData []blogs.BlogInfo) (string, error) {
	var buffer bytes.Buffer
	if err := blogEntriesTemplate.Execute(&buffer, blogsData); err != nil {
		return "", fmt.Errorf("error parsing blog entries: %w", err)
	}
	return buffer.String(), nil
}

func NewBlogsAPI(repo *blogs.Repository) *BlogsAPI {
	return &BlogsAPI{repo: repo}
}

func (a *BlogsAPI) Read(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	var kind blogs.Kind
	var items []blogs.BlogInfo
	var err error

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
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Default to HTML for text/html, */* or empty Accept header
	if accept == "" || accept == "*/*" || strings.Contains(accept, "text/html") {
		htmlContent, err := blogsDataToCards(items)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset: utf-8")
		io.WriteString(w, htmlContent)
		return
	}

	http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
}
