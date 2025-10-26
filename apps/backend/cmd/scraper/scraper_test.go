package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nesco/techblogs/backend/internal/blogs"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "https with absolute path",
			baseURL:  "https://example.com",
			path:     "/blog/post",
			expected: "https://example.com/blog/post",
		},
		{
			name:     "https with relative path",
			baseURL:  "https://example.com",
			path:     "post",
			expected: "https://example.com/post",
		},
		{
			name:     "http with absolute path",
			baseURL:  "http://example.com",
			path:     "/blog/post",
			expected: "http://example.com/blog/post",
		},
		{
			name:     "base with trailing slash",
			baseURL:  "https://example.com/",
			path:     "/post",
			expected: "https://example.com/post",
		},
		{
			name:     "base with path should strip it",
			baseURL:  "https://example.com/blog/archives",
			path:     "/post-1",
			expected: "https://example.com/post-1",
		},
		{
			name:     "no scheme defaults to https",
			baseURL:  "example.com",
			path:     "/post",
			expected: "https://example.com/post",
		},
		{
			name:     "path without leading slash",
			baseURL:  "https://example.com",
			path:     "posts/latest",
			expected: "https://example.com/posts/latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeURL(tt.baseURL, tt.path)
			if got != tt.expected {
				t.Errorf("normalizeURL(%q, %q) = %q, want %q", tt.baseURL, tt.path, got, tt.expected)
			}
		})
	}
}

func TestScrapeBlog(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		hrefSelector  string
		nameSelector  string
		expectedName  string
		expectedHref  string
		expectError   bool
		errorContains string
	}{
		{
			name: "basic article with absolute href",
			html: `<html><body>
				<a href="https://example.com/post-1" class="article-link">My First Post</a>
			</body></html>`,
			hrefSelector: "a.article-link",
			nameSelector: "a.article-link",
			expectedName: "My First Post",
			expectedHref: "https://example.com/post-1",
			expectError:  false,
		},
		{
			name: "article with relative href",
			html: `<html><body>
				<a href="/posts/latest" class="link">Latest Article</a>
			</body></html>`,
			hrefSelector: "a.link",
			nameSelector: "a.link",
			expectedName: "Latest Article",
			expectError:  false,
		},
		{
			name: "article with nested text elements",
			html: `<html><body>
				<a href="/post" class="card">
					<h2 class="title">
						Understanding Go Testing
					</h2>
				</a>
			</body></html>`,
			hrefSelector: "a.card",
			nameSelector: "h2.title",
			expectedName: "Understanding Go Testing",
			expectError:  false,
		},
		{
			name: "article with excessive whitespace",
			html: `<html><body>
				<a href="/post" class="link">

					  Multi
					  Line
					  Title

				</a>
			</body></html>`,
			hrefSelector: "a.link",
			nameSelector: "a.link",
			expectedName: "Multi Line Title",
			expectError:  false,
		},
		{
			name: "article with aria-label fallback",
			html: `<html><body>
				<a href="/post" class="link" aria-label="Accessible Title"></a>
			</body></html>`,
			hrefSelector: "a.link",
			nameSelector: "a.link",
			expectedName: "Accessible Title",
			expectError:  false,
		},
		{
			name:          "no article found with href selector",
			html:          `<html><body><div>No articles here</div></body></html>`,
			hrefSelector:  "a.article",
			nameSelector:  "a.article",
			expectError:   true,
			errorContains: "no articles found with href selector",
		},
		{
			name: "article without href attribute",
			html: `<html><body>
				<a class="link">Article Title</a>
			</body></html>`,
			hrefSelector:  "a.link",
			nameSelector:  "a.link",
			expectError:   true,
			errorContains: "article href not found",
		},
		{
			name: "no article found with name selector",
			html: `<html><body>
				<a href="/post" class="link"></a>
			</body></html>`,
			hrefSelector:  "a.link",
			nameSelector:  "h2.title",
			expectError:   true,
			errorContains: "no articles found with name selector",
		},
		{
			name: "empty article name",
			html: `<html><body>
				<a href="/post" class="link"></a>
			</body></html>`,
			hrefSelector:  "a.link",
			nameSelector:  "a.link",
			expectError:   true,
			errorContains: "article name is empty",
		},
		{
			name:         "multiple articles selects first",
			html: `<html><body>
				<a href="/post-1" class="article">First Post</a>
				<a href="/post-2" class="article">Second Post</a>
			</body></html>`,
			hrefSelector: "a.article",
			nameSelector: "a.article",
			expectedName: "First Post",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			config := blogs.BlogConfig{
				BlogName:            "Test Blog",
				BlogHref:            server.URL,
				ArticleHrefSelector: tt.hrefSelector,
				ArticleNameSelector: tt.nameSelector,
			}

			// Call scrapeBlog
			name, href, err := scrapeBlog(config)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorContains)
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if name != tt.expectedName {
				t.Errorf("name = %q, want %q", name, tt.expectedName)
			}

			// For relative URLs, expect them to be normalized with server URL
			if tt.expectedHref != "" && !strings.HasPrefix(tt.expectedHref, "http") {
				tt.expectedHref = server.URL + tt.expectedHref
			}

			if tt.expectedHref != "" && href != tt.expectedHref {
				t.Errorf("href = %q, want %q", href, tt.expectedHref)
			}
		})
	}
}

func TestScrapeBlog_HTTPErrors(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		errorContains string
	}{
		{
			name:          "404 not found",
			statusCode:    http.StatusNotFound,
			errorContains: "bad status code: 404",
		},
		{
			name:          "500 internal server error",
			statusCode:    http.StatusInternalServerError,
			errorContains: "bad status code: 500",
		},
		{
			name:          "403 forbidden",
			statusCode:    http.StatusForbidden,
			errorContains: "bad status code: 403",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			config := blogs.BlogConfig{
				BlogName:            "Test Blog",
				BlogHref:            server.URL,
				ArticleHrefSelector: "a.article",
				ArticleNameSelector: "a.article",
			}

			_, _, err := scrapeBlog(config)
			if err == nil {
				t.Errorf("expected error, got nil")
			} else if !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
			}
		})
	}
}

func TestScrapeBlog_EmptySelector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body></body></html>`))
	}))
	defer server.Close()

	config := blogs.BlogConfig{
		BlogName:            "Test Blog",
		BlogHref:            server.URL,
		ArticleHrefSelector: "",
		ArticleNameSelector: "",
	}

	name, href, err := scrapeBlog(config)
	if err != nil {
		t.Fatalf("expected no error for empty selector, got: %v", err)
	}

	if name != "" || href != "" {
		t.Errorf("expected empty results, got name=%q href=%q", name, href)
	}
}

func TestScrapeBlog_MalformedHTML(t *testing.T) {
	// Test that goquery handles malformed HTML gracefully
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body><a href="/post" class="link">Unclosed tag`))
	}))
	defer server.Close()

	config := blogs.BlogConfig{
		BlogName:            "Test Blog",
		BlogHref:            server.URL,
		ArticleHrefSelector: "a.link",
		ArticleNameSelector: "a.link",
	}

	name, href, err := scrapeBlog(config)
	if err != nil {
		t.Fatalf("goquery should handle malformed HTML gracefully, got error: %v", err)
	}

	if name != "Unclosed tag" {
		t.Errorf("name = %q, want %q", name, "Unclosed tag")
	}

	if href != server.URL+"/post" {
		t.Errorf("href = %q, want %q", href, server.URL+"/post")
	}
}
