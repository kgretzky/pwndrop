package core

import (
	"context"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

type CertDb struct {
	AutocertMgr autocert.Manager
}

func NewCertDb(cache_dir string) (*CertDb, error) {
	cdb := &CertDb{
		AutocertMgr: autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(filepath.Join(cache_dir, "autocert")),
		},
	}
	cdb.AutocertMgr.HostPolicy = cdb.hostPolicy
	return cdb, nil
}

func (cdb *CertDb) SetManagedHostnames(hosts ...string) {
	cdb.AutocertMgr.HostPolicy = autocert.HostWhitelist(hosts...)
}

func (cdb *CertDb) hostPolicy(ctx context.Context, host string) error {
	// accept all hosts for TLS certificate retrieval
	return nil
}
