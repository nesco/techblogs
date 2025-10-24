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
		<p>Latest: <a href="{{ .LatestArticleHref }}">{{ .LatestArticleName }}</a></p>
	</article>
	{{- end -}}
`))

type BlogsAPI struct {
	blogsRepository []blogs.BlogInfo
}

func blogsDataToCards(blogsData []blogs.BlogInfo) (string, error) {
	var buffer bytes.Buffer
	if err := blogEntriesTemplate.Execute(&buffer, blogsData); err != nil {
		return "", fmt.Errorf("error parsing blog entries: %w", err)
	}
	return buffer.String(), nil
}

func NewBlogsAPI() *BlogsAPI {
	orgData := []blogs.BlogInfo{
		{
			BlogHref:          "https://stripe.com/blog",
			BlogName:          "Stripe",
			LatestArticleHref: "https://stripe.com/blog/introducing-stablecoin-payments-for-subscriptions",
			LatestArticleName: "Introducing stablecoin payments for subscriptions",
			Kind:              blogs.Organization,
		},
		{
			BlogHref:          "https://www.datadoghq.com/blog/",
			BlogName:          "Datadoghq",
			LatestArticleHref: "https://www.datadoghq.com/state-of-cloud-security/",
			LatestArticleName: "State of Cloud Security",
			Kind:              blogs.Organization,
		},
	}

	peopleData := []blogs.BlogInfo{
		{
			BlogHref:          "https://buttondown.com/hillelwayne/archive/",
			BlogName:          "Hillel Wayne",
			LatestArticleHref: "https://buttondown.com/hillelwayne/archive/modal-editing-is-a-weird-historical-contingency/",
			LatestArticleName: "Modal editing is a weird historical contingency we have through sheer happenstance",
			Kind:              blogs.Person,
		},
		{
			BlogHref:          "https://blog.samaltman.com/",
			BlogName:          "Sam Altman",
			LatestArticleHref: "https://blog.samaltman.com/sora-update-number-1",
			LatestArticleName: "Sora update #1",
			Kind:              blogs.Person,
		},
	}

	blogsRepository := append(orgData, peopleData...)
	return &BlogsAPI{blogsRepository}

}

func (a *BlogsAPI) Read(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	var kind blogs.Kind

	if collection != "" {
		var ok bool
		kind, ok = blogs.KindByCollection[collection]
		if !ok {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}
	}

	var items []blogs.BlogInfo
	for _, blog := range a.blogsRepository {
		if blog.Kind == kind || kind == "" {
			items = append(items, blog)
		}
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
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, htmlContent)
		return
	}

	// Unsupported media type
	http.Error(w, "Not Acceptable", http.StatusNotAcceptable)

}
