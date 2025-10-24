package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

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
	startTime := time.Now()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Infow("Techblogs backend started")

	// Setting-up the HTTP Server
	mux := http.NewServeMux()

	registerRoutes(mux, startTime)

	addr := os.Getenv("LISTEN_ADDR")
	ln, err := getListener(addr)

	if err != nil {
		sugar.Fatal(err)
	}

	http.Serve(ln, mux)

}
