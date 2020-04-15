package api

import (
	"net/http"

	"github.com/kgretzky/pwndrop/config"
)

func VersionOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
}

func VersionGetHandler(w http.ResponseWriter, r *http.Request) {
	type VersionResponse struct {
		Version string `json:"version"`
	}

	resp := &VersionResponse{}
	resp.Version = config.Version

	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}
