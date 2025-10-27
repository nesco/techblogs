// Package home renders the home page of techblo.gs
package home

import (
	"bytes"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type HomeHandler struct {
	logger zap.SugaredLogger
}

func (a *HomeHandler) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buffer bytes.Buffer
	if err := homePageTemplate.Execute(&buffer, nil); err != nil {
		a.logger.Errorf("error parsing main page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err := io.WriteString(w, buffer.String())
	if err != nil {
		a.logger.Errorw("error encoding home page response", "error", err, "client", r.Header.Get("X-Real-IP"))
	}
}
