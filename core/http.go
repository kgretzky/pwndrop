package core

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kgretzky/pwndrop/log"
)

type Http struct {
	srv *Server
}

func NewHttp(srv *Server) (*Http, error) {
	s := &Http{
		srv: srv,
	}
	return s, nil
}

func (s *Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data_dir := Cfg.GetDataDir()
	if r.Method == "GET" {
		f, status, err := s.srv.GetFile(r.URL.Path)
		if err != nil {
			w.WriteHeader(status)
			log.Error("http: get: %s: %s", r.URL.Path, err)
			return
		}

		if f.RedirectPath != "" && f.RedirectPath != r.URL.Path && !f.IsPaused {
			http.Redirect(w, r, f.RedirectPath, http.StatusFound)
		} else {
			mime_type := f.MimeType
			if f.IsPaused {
				mime_type = f.SubMimeType
			}
			fpath := filepath.Join(data_dir, "files", f.Filename)
			fo, err := os.Open(fpath)
			//data, err := ioutil.ReadFile(fpath)
			if err != nil {
				log.Error("http: file: %s: %s", f.Filename, err)
				return
			}
			defer fo.Close()

			w.Header().Set("Content-Type", mime_type)
			w.WriteHeader(200)
			io.Copy(w, fo)
		}
		return
	}
	w.WriteHeader(404)
}
