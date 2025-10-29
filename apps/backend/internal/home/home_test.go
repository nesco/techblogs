package home

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nesco/techblogs/backend/internal/blogs"
	"go.uber.org/zap"
)

func TestGetHome(t *testing.T) {
	// Create in-memory database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE blog_cache (
			id INTEGER PRIMARY KEY,
			blog_name TEXT NOT NULL UNIQUE,
			blog_href TEXT NOT NULL,
			latest_article_name TEXT,
			latest_article_href TEXT,
			kind TEXT NOT NULL CHECK (kind IN ('organization', 'individual')),
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO blog_cache (blog_name, blog_href, latest_article_name, latest_article_href, kind)
		VALUES ('Test Blog', 'https://example.com', 'Test Article', 'https://example.com/article', 'individual')
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := blogs.NewRepository(db)
	api := HomeHandler{Logger: *logger, Repo: repo}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	api.Read(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected content-type 'text/html; charset=utf-8', got %s", res.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read the response body %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty response body")
	}

	doctype := "<!doctype html>"
	bodyStart := strings.ToLower(string(body)[:len(doctype)])
	if bodyStart != doctype {
		t.Errorf("expected body to start with %q, got %q", doctype, bodyStart)
	}
}
