package home

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
)

var api = HomeAPI{zap.SugaredLogger{}}

func TestGetHome(t *testing.T) {
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
