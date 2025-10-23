package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

var startTime = time.Now()

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().UTC().Format(time.RFC3339),
		"uptime_s":  time.Since(startTime).Seconds(),
	})
}

func organisationEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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

	ln, err := getListener("")

	if err != nil {
		sugar.Fatal(err)
	}

	http.Serve(ln, mux)

}
