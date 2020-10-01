package core

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme"

	"github.com/kgretzky/pwndrop/api"
	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
)

const (
	API_PATH = "api/v1"
)

type Server struct {
	srv       *http.Server
	listenTLS net.Listener
	listen    net.Listener
	wdav      *WebDav
	http      *Http
	cdb       *CertDb
	ns        *Nameserver
	r         *mux.Router
	blacklist map[string]*BlacklistItem
	bl_mtx    sync.Mutex
}

func NewServer(host string, port_plain int, port_tls int, enable_letsencrypt bool, enable_dns bool, ch_exit *chan bool) (*Server, error) {
	var err error
	s := &Server{
		blacklist: make(map[string]*BlacklistItem),
		bl_mtx:    sync.Mutex{},
	}

	hostname := fmt.Sprintf("%s:%d", host, port_plain)
	hostname_tls := fmt.Sprintf("%s:%d", host, port_tls)

	s.cdb, err = NewCertDb(Cfg.GetDataDir())
	if err != nil {
		return nil, err
	}

	cert, err := LoadTLSCertificate(filepath.Join(Cfg.GetDataDir(), "public.crt"), filepath.Join(Cfg.GetDataDir(), "private.key"))
	if err != nil {
		log.Warning("certificate: %s", err)
		cert, err = GenerateTLSCertificate(host)
		if err != nil {
			return nil, err
		}
		log.Info("generated self-signed certificate")
	} else {
		log.Info("using TLS certificate from data directory")
		enable_letsencrypt = false
	}

	tls_cfg := &tls.Config{}
	tls_cfg.Certificates = append(tls_cfg.Certificates, *cert)
	if enable_letsencrypt {
		log.Info("autocert: enabled")
		tls_cfg.GetCertificate = s.cdb.AutocertMgr.GetCertificate
		tls_cfg.NextProtos = []string{
			"h2", "http/1.1", // enable HTTP/2
			acme.ALPNProto, // enable tls-alpn ACME challenges
		}
	} else {
		log.Info("autocert: disabled")
	}

	// set up modern cipher suites
	/*
		tls_cfg.MinVersion = tls.VersionTLS12
		tls_cfg.CipherSuites = []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		}*/

	s.wdav, err = NewWebDav(s)
	if err != nil {
		return nil, err
	}
	s.http, err = NewHttp(s)
	if err != nil {
		return nil, err
	}

	s.setupRouter()

	s.srv = &http.Server{
		Handler:      http.Handler(s),
		Addr:         hostname,
		WriteTimeout: 0,
		ReadTimeout:  0,
		IdleTimeout:  5 * time.Second,
		TLSConfig:    tls_cfg,
	}

	s.listenTLS, err = tls.Listen("tcp", hostname_tls, tls_cfg)
	if err != nil {
		return nil, err
	}
	s.listen, err = net.Listen("tcp", hostname)
	if err != nil {
		return nil, err
	}

	log.Info("starting HTTP/WebDAV server at %s", hostname)
	log.Info("starting HTTPS server at %s", hostname_tls)

	if enable_dns {
		s.ns, err = NewNameserver(ch_exit)
		if err != nil {
			return nil, err
		}
	}

	go func() {
		err := s.srv.Serve(s.listen)
		if err != nil {
			log.Fatal("failed to start HTTP/WebDAV server at %s", hostname)
			*ch_exit <- false
		}
	}()

	go func() {
		err := s.srv.Serve(s.listenTLS)
		if err != nil {
			log.Fatal("failed to start HTTPS server at %s", hostname_tls)
			*ch_exit <- false
		}
	}()

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("%s %s", r.Method, r.URL.Path)

	from_ip := r.RemoteAddr
	if strings.Contains(from_ip, ":") {
		from_ip = strings.Split(from_ip, ":")[0]
	}

	if s.isBlacklisted(from_ip) {
		err := s.killConnection(w, -1)
		if err != nil {
			log.Error("http: %s (%s)", err, from_ip)
			w.Header().Set("Connection", "close")
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(500)
		}
		return
	}

	if !s.isWebDavRequest(r) {

		cookie_name := Cfg.GetCookieName()
		cookie_token := Cfg.GetCookieToken()

		if r.URL.Path == Cfg.GetSecretPath() {
			ck := &http.Cookie{
				Domain:   "",
				Path:     "/",
				Expires:  time.Now().AddDate(0, 3, 0),
				HttpOnly: true,
				Name:     cookie_name,
				Value:    cookie_token,
			}
			http.SetCookie(w, ck)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if !s.FileExists(r.URL.Path) {
			if ck, err := r.Cookie(cookie_name); err == nil {
				if ck.Value == cookie_token {
					s.r.ServeHTTP(w, r)
					return
				}
			}

			s.addBlacklistHit(from_ip)
			if len(Cfg.GetRedirectUrl()) > 0 {
				http.Redirect(w, r, Cfg.GetRedirectUrl(), http.StatusFound)
				return
			}
		}

		// http
		s.http.ServeHTTP(w, r)
	} else {
		// webdav
		s.wdav.Handler().ServeHTTP(w, r)
	}
}

func (s *Server) isWebDavRequest(r *http.Request) bool {
	ua := r.Header.Get("user-agent")
	if strings.Index(ua, "WebDAV") >= 0 || strings.Index(ua, "DavClnt") >= 0 {
		return true
	}
	if r.Header.Get("translate") == "f" {
		return true
	}
	return false
}

func (s *Server) setupRouter() {
	admin_path := "/"
	s.r = mux.NewRouter()
	sr := s.r.PathPrefix(admin_path + API_PATH).Subrouter()
	sr.HandleFunc("/auth", api.AuthOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/auth", api.AuthCheckHandler).Methods("GET")
	sr.HandleFunc("/server_info", api.ServerInfoOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/server_info", api.ServerInfoGetHandler).Methods("GET")
	sr.HandleFunc("/version", api.VersionOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/version", api.VersionGetHandler).Methods("GET")
	sr.HandleFunc("/login", api.AuthOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/login", api.LoginUserHandler).Methods("POST")
	sr.HandleFunc("/logout", api.LogoutUserHandler).Methods("GET")
	sr.HandleFunc("/clear_secret", api.ClearSecretSessionHandler).Methods("GET")
	sr.HandleFunc("/create_account", api.AuthOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/create_account", api.CreateUserHandler).Methods("POST")
	sr.HandleFunc("/config", api.ConfigOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/config", api.ConfigGetHandler).Methods("GET")
	sr.HandleFunc("/config", api.ConfigUpdateHandler).Methods("POST")
	sr.HandleFunc("/files", api.FileOptionsHandler).Methods("OPTIONS")
	sr.HandleFunc("/files", api.FileListHandler).Methods("GET")
	sr.HandleFunc("/files", api.FileCreateHandler).Methods("POST")
	sr.HandleFunc("/files/{id}", api.FileDeleteHandler).Methods("DELETE")
	sr.HandleFunc("/files/{id}", api.FileUpdateHandler).Methods("PUT")
	sr.HandleFunc("/files/{id}/sub", api.SubFileCreateHandler).Methods("POST")
	sr.HandleFunc("/files/{id}/sub/{sub_id}", api.SubFileDeleteHandler).Methods("DELETE")
	sr.HandleFunc("/files/{id}/enable", api.FileEnableHandler).Methods("GET")
	sr.HandleFunc("/files/{id}/disable", api.FileDisableHandler).Methods("GET")
	sr.HandleFunc("/files/{id}/pause", api.FilePauseHandler).Methods("GET")
	sr.HandleFunc("/files/{id}/unpause", api.FileUnpauseHandler).Methods("GET")
	s.r.PathPrefix(fmt.Sprintf("%s", admin_path)).Handler(http.StripPrefix(fmt.Sprintf("%s", admin_path), http.FileServer(http.Dir(Cfg.GetAdminDir()))))
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Debug("%s %s", r.Method, r.URL.Path)
}

func (s *Server) GetFile(url string) (*storage.DbFile, int, error) {
	is_redirect := false
	f, err := storage.FileGetByUrl(url)
	if err != nil {
		f, err = storage.FileGetByRedirectUrl(url)
		if err != nil {
			return nil, 404, err
		}
		is_redirect = true
	}
	if !f.IsEnabled {
		return nil, 404, fmt.Errorf("file is disabled")
	}
	if f.IsPaused {
		if f.RedirectPath != "" && is_redirect {
			return nil, 404, fmt.Errorf("can't access facade via redirect while paused")
		} else if f.RefSubFile > 0 {
			sf, err := storage.SubFileGet(f.RefSubFile)
			if err != nil {
				return nil, 404, fmt.Errorf("facade file not found")
			}
			f.Filename = sf.Filename
			f.FileSize = sf.FileSize
		} else {
			return nil, 404, fmt.Errorf("facade file not set")
		}
	}
	return f, 200, nil
}

func (s *Server) FileExists(url string) bool {
	_, err := storage.FileGetByUrl(url)
	if err != nil {
		_, err = storage.FileGetByRedirectUrl(url)
		if err != nil {
			return false
		}
	}
	return true
}

func (s *Server) killConnection(w http.ResponseWriter, status int) error {
	if status > 0 {
		w.Header().Set("Connection", "close")
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(status)
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		return fmt.Errorf("connection hijacking not supported")
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (s *Server) isBlacklisted(ip_addr string) bool {
	s.bl_mtx.Lock()
	defer s.bl_mtx.Unlock()

	ret := false
	if bl, ok := s.blacklist[ip_addr]; ok {
		if bl.hits >= BLACKLIST_HITS_LIMIT {
			if time.Now().Before(bl.last_hit.Add(BLACKLIST_JAIL_TIME_SECS * time.Second)) {
				ret = true
			} else {
				delete(s.blacklist, ip_addr)
				return false
			}
		}
		bl.last_hit = time.Now()
	}
	return ret
}

func (s *Server) addBlacklistHit(ip_addr string) {
	s.bl_mtx.Lock()
	defer s.bl_mtx.Unlock()

	if bl, ok := s.blacklist[ip_addr]; ok {
		bl.hits += 1
	} else {
		bl := &BlacklistItem{
			hits:     1,
			last_hit: time.Now(),
		}
		s.blacklist[ip_addr] = bl
	}
}
