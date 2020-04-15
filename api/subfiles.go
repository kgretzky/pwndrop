package api

import (
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

func SubFileCreateHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)

	data_dir := Cfg.GetDataDir()
	user_id := 1

	file, fhead, err := r.FormFile("file")
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}
	defer file.Close()

	fid, err := strconv.Atoi(vars["id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	parent_file, err := storage.FileGet(fid)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_FILE_NOT_FOUND, nil)
		return
	}

	name := fhead.Filename
	fname := utils.GenRandomHash()

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

	o := &storage.DbSubFile{
		Fid:        fid,
		Uid:        user_id,
		Name:       name,
		Filename:   fname,
		FileSize:   fi.Size(),
		CreateTime: time.Now().Unix(),
	}

	f, err := storage.SubFileCreate(o)
	if err != nil {
		os.Remove(save_path)
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	parent_file.SubName = f.Name
	parent_file.RefSubFile = f.ID
	log.Debug("ref_sub_file: %d", parent_file.RefSubFile)
	_, err = storage.FileUpdate(fid, parent_file)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_SAVE_FAILED, nil)
		return
	}

	DumpResponse(w, "ok", http.StatusOK, 0, f)
}

func SubFileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	vars := mux.Vars(r)
	sub_id, err := strconv.Atoi(vars["sub_id"])
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	err = DeleteSubFile(sub_id)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	files, err := storage.FileList()
	if err == nil {
		for _, f := range files {
			if f.RefSubFile == sub_id {
				storage.FilePause(f.ID, false)
			}
		}
	}

	DumpResponse(w, "ok", http.StatusOK, 0, nil)
}

func DeleteSubFile(sub_id int) error {
	data_dir := Cfg.GetDataDir()
	f, err := storage.SubFileGet(sub_id)
	if err != nil {
		return err
	}

	_, err = storage.FileResetSubFile(f.Fid)
	if err != nil {
		return err
	}

	err = storage.SubFileDelete(sub_id)
	if err != nil {
		return err
	}
	save_path := filepath.Join(data_dir, "files", f.Filename)
	os.Remove(save_path)
	return nil
}
