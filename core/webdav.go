package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/webdav"

	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
)

type WebDav struct {
	srv *Server
	h   http.Handler
}

func NewWebDav(srv *Server) (*WebDav, error) {
	s := &WebDav{
		srv: srv,
	}

	fs := &WebDavFS{
		srv: srv,
	}

	s.h = &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Debug("WEBDAV [%s]: %s, ERROR: %s", r.Method, r.URL, err)
			} else {
				log.Debug("WEBDAV [%s]: %s", r.Method, r.URL)
			}
		},
	}

	return s, nil
}

func (s *WebDav) Handler() http.Handler {
	return s.h
}

// ----

type WebDavFS struct {
	srv *Server
}

func (fs *WebDavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return fmt.Errorf("mkdir: not supported")
}

func (fs *WebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	//return nil, fmt.Errorf("openfile: not supported")

	log.Debug("openfile: %s %d %08x", name, flag, perm)

	is_dir := storage.FileDirExists(name)

	var fsize int64 = 0
	f, _, err := fs.srv.GetFile(name)
	if err != nil {
		if !is_dir {
			log.Error("webdav: %s", err)
			return nil, err
		}
	} else {
		fsize = f.FileSize
	}

	fi := &WebDavFileInfo{
		name:    name,
		size:    fsize,
		isDir:   is_dir,
		modTime: time.Now(),
	}

	wf := &WebDavFile{
		fi: fi,
	}

	if !is_dir {
		data_dir := Cfg.GetDataDir()
		fpath := filepath.Join(data_dir, "files", f.Filename)
		wf.fh, err = os.OpenFile(fpath, flag, perm)
		if err != nil {
			log.Error("webdav: %s", err)
			return nil, err
		}
	}
	return wf, nil
}

func (fs *WebDavFS) RemoveAll(ctx context.Context, name string) error {
	return fmt.Errorf("removeall: not supported")
}

func (fs *WebDavFS) Rename(ctx context.Context, oldName, newName string) error {
	return fmt.Errorf("rename: not supported")
}

func (fs *WebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Debug("webdav: stat('%s')", name)
	is_dir := false
	if name == "" {
		return nil, fmt.Errorf("invalid name")
	}

	if name[len(name)-1] == '/' {
		is_dir = true
	}

	fi := &WebDavFileInfo{
		name:    name,
		isDir:   is_dir,
		modTime: time.Now(),
	}
	is_dir = storage.FileDirExists(name)

	if !is_dir {
		f, _, err := fs.srv.GetFile(name)
		if err != nil {
			log.Error("webdav: %s", err)
			return nil, err
		}
		fi.size = f.FileSize
	}

	return fi, nil
}

// ----

type WebDavFileInfo struct {
	os.FileInfo

	name    string
	size    int64
	isDir   bool
	modTime time.Time
}

func (fi *WebDavFileInfo) Name() string {
	return fi.name
}

func (fi *WebDavFileInfo) Size() int64 {
	return fi.size
}

func (fi *WebDavFileInfo) Mode() os.FileMode {
	var ret os.FileMode = 0644
	if fi.isDir {
		ret |= os.ModeDir
	}
	return ret
}

func (fi *WebDavFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *WebDavFileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *WebDavFileInfo) Sys() interface{} {
	return nil
}

// ----

type WebDavFile struct {
	webdav.File

	fi os.FileInfo
	fh *os.File
}

func (f *WebDavFile) Readdir(count int) ([]os.FileInfo, error) {
	log.Debug("readdir: %d", count)
	return []os.FileInfo{}, nil
}

func (f *WebDavFile) Stat() (os.FileInfo, error) {
	log.Debug("webdavfile: stat()")
	return f.fi, nil
}

func (f *WebDavFile) Close() error {
	log.Debug("webdavfile: close()")
	return f.fh.Close()
}

func (f *WebDavFile) Read(p []byte) (n int, err error) {
	log.Debug("webdavfile: read()")
	return f.fh.Read(p)
}

func (f *WebDavFile) Seek(offset int64, whence int) (int64, error) {
	log.Debug("webdavfile: seek(%d, %d)", offset, whence)
	return f.fh.Seek(offset, whence)
}

func (f *WebDavFile) Write(p []byte) (n int, err error) {
	log.Debug("webdavfile: write()")
	return 0, fmt.Errorf("write not supported")
}
