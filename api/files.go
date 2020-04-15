package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
	"github.com/kgretzky/pwndrop/utils"
)

func FileOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
}

func FileCreateHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	data_dir := Cfg.GetDataDir()
	user_id := 1

	file, fhead, err := r.FormFile("file")
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	defer file.Close()

	name := fhead.Filename
	fname := utils.GenRandomHash()
	url_path := "/" + utils.GenRandomString(8) + "/" + name // TODO: make sure the generated folder is unique
	mime_type := fhead.Header.Get("content-type")           //r.Header.Get("content-type")
	if mime_type == "" {
		mime_type = "application/octet-stream"
	}
	log.Debug("upload: " + mime_type)

	os.Mkdir(filepath.Join(data_dir, "files"), 0700)
	save_path := filepath.Join(data_dir, "files", fname)
	if err := SaveUploadedFile(file, fhead, save_path); err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_SAVE_FAILED, nil)
		return
	}

	var fi os.FileInfo
	if fi, err = os.Stat(save_path); err != nil {
		os.Remove(save_path)
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_NOT_FOUND, nil)
		return
	}

	o := &storage.DbFile{
		Uid:          user_id,
		Name:         name,
		Filename:     fname,
		FileSize:     fi.Size(),
		UrlPath:      url_path,
		RedirectPath: "",
		MimeType:     mime_type,
		SubMimeType:  mime_type,
		OrigMimeType: mime_type,
		CreateTime:   time.Now().Unix(),
		IsEnabled:    true,
		IsPaused:     false,
		RefSubFile:   0,
	}

	f, err := storage.FileCreate(o)
	if err != nil {
		os.Remove(save_path)
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func FileListHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	files, err := storage.FileList()
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	type JsonFile struct {
		storage.DbFile
		SubFile *storage.DbSubFile `json:"sub_file"`
	}
	type Response struct {
		Uploads []*JsonFile `json:"uploads"`
	}
	resp := &Response{}

	for _, f := range files {
		jo := &JsonFile{
			DbFile: f,
		}
		if f.RefSubFile > 0 {
			subf, err := storage.SubFileGet(f.RefSubFile)
			if err == nil {
				jo.SubFile = subf
			}
		}
		resp.Uploads = append(resp.Uploads, jo)
	}

	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	data_dir := Cfg.GetDataDir()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	f, err := storage.FileGet(id)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusNotFound, API_ERROR_FILE_NOT_FOUND, nil)
		return
	}

	if f.RefSubFile > 0 {
		err = DeleteSubFile(f.RefSubFile)
		if err != nil {
			DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
			return
		}
	}

	err = storage.FileDelete(id)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	save_path := filepath.Join(data_dir, "files", f.Filename)
	os.Remove(save_path)

	DumpResponse(w, "ok", http.StatusOK, 0, nil)
}

func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	file := storage.DbFile{}
	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	if file.UrlPath[0] != '/' {
		file.UrlPath = "/" + file.UrlPath
	}
	if len(file.RedirectPath) > 0 && file.RedirectPath[0] != '/' {
		file.RedirectPath = "/" + file.RedirectPath
	}

	f, err := storage.FileUpdate(id, &file)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func FileEnableHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	f, err := storage.FileEnable(id, true)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func FileDisableHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	f, err := storage.FileEnable(id, false)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	log.Debug("%v", f.IsEnabled)
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func FilePauseHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	f, err := storage.FilePause(id, true)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func FileUnpauseHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	f, err := storage.FilePause(id, false)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, f)
}
