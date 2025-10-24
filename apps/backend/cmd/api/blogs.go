package main

import (
	"bytes"
	"fmt"
	"github.com/nesco/techblogs/backend/internal/blogs"
	"html/template"
	"io"
	"net/http"
)

type BlogsAPI struct {
	blogsRepository []blogs.BlogInfo
}

func blogsDataToCards(blogsData []blogs.BlogInfo) (string, error) {
	const blogEntriesCardTemplate = `
		{{- range . -}}
		<article class="card">
		 <h3><a href="{{ .BlogHref }}">{{ .BlogName }}</a></h3>
			<p>Latest: <a href="{{ .LatestArticleHref }}">{{ .LatestArticleName }}</a></p>
		</article>
		{{- end -}}
		`
	var buffer bytes.Buffer
	templateParsed := template.Must(template.New("BlogEntries").Parse(blogEntriesCardTemplate))
	if err := templateParsed.Execute(&buffer, blogsData); err != nil {
		return "", fmt.Errorf("Error parsing blog entries: %w", err)
	}

	return buffer.String(), nil
}

func NewBlogsAPI() *BlogsAPI {
	var orgData = []blogs.BlogInfo{
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

	var peopleData = []blogs.BlogInfo{
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
	kind := blogs.KindByCollection[r.PathValue("collection")]
	var items []blogs.BlogInfo
	for _, blog := range a.blogsRepository {
		if blog.Kind == kind || kind == "" {
			items = append(items, blog)
		}
	}

	var htmlContent string
	var err error
	if htmlContent, err = blogsDataToCards(items); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, htmlContent)

}
