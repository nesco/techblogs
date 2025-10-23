package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

type BlogInfo struct {
	BlogHref          string `json:"blogHref"`
	BlogName          string `json:"blogName"`
	LatestArticleHref string `json:"latestArticleHref"`
	LatestArticleName string `json:"latestArticleName"`
}

var orgData = []BlogInfo{{"https://stripe.com/blog", "Stripe", "https://stripe.com/blog/introducing-stablecoin-payments-for-subscriptions", "Introducing stablecoin payments for subscriptions"}}
var startTime = time.Now()

func blogsDataToCards(blogsData []BlogInfo) (string, error) {
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
		return "", fmt.Errorf("Error parsing blog entries card template: %w", err)
	}

	return buffer.String(), nil

}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"timestamp": time.Now().UTC().UTC().Format(time.RFC3339),
			"uptime_s":  time.Since(startTime).Seconds(),
		})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func organisationEndpoint(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		var htmlContent string
		var err error
		if htmlContent, err = blogsDataToCards(orgData); err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, htmlContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getListener(addr string) (net.Listener, error) {

	if addr == "" {
		addr = "127.0.0.1:5011"
	}

	addr = strings.TrimPrefix(addr, "tcp://")

	if sock, found := strings.CutPrefix(addr, "unix:"); found {
		_ = os.Remove(sock)

		if err := os.MkdirAll(filepath.Dir(sock), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", sock, err)
		}

		ln, err := net.Listen("unix", sock)
		if err != nil {
			return nil, fmt.Errorf("listen unix %s: %w", sock, err)
		}

		// Make sure nginx/app user can read/write
		if err := os.Chmod(sock, 0o660); err != nil {
			_ = ln.Close()
			return nil, fmt.Errorf("chmod %s: %w", sock, err)
		}

		return ln, nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen tcp %s: %w", addr, err)
	}

	return ln, nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Infow("Techblogs backend started")

	// Setting-up the HTTP Server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthEndpoint)
	mux.HandleFunc("/organisations", organisationEndpoint)

	addr := os.Getenv("LISTEN_ADDR")
	ln, err := getListener(addr)

	if err != nil {
		sugar.Fatal(err)
	}

	http.Serve(ln, mux)

}
