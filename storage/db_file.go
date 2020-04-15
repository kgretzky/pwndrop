package storage

import (
	"github.com/kgretzky/pwndrop/log"
)

type DbFile struct {
	ID           int    `json:"id" storm:"id,increment"`
	Uid          int    `json:"uid" storm:"index"`
	Name         string `json:"name"`
	Filename     string `json:"fname"`
	FileSize     int64  `json:"fsize"`
	UrlPath      string `json:"url_path" storm:"unique"`
	MimeType     string `json:"mime_type"`
	OrigMimeType string `json:"orig_mime_type"`
	CreateTime   int64  `json:"create_time" storm:"index"`
	IsEnabled    bool   `json:"is_enabled"`
	IsPaused     bool   `json:"is_paused"`
	RedirectPath string `json:"redirect_path" storm:"unique"`
	SubName      string `json:"sub_name"`
	SubMimeType  string `json:"sub_mime_type"`
	RefSubFile   int    `json:"ref_sub_file"`
}

func FileCreate(o *DbFile) (*DbFile, error) {
	err := db.Save(o)
	if err != nil {
		return nil, err
	}
	log.Debug("file id: %d", o.ID)
	return o, nil
}

func FileList() ([]DbFile, error) {
	var dbos []DbFile

	err := db.All(&dbos)
	if err != nil {
		return nil, err
	}
	/*
		for _, dbo := range dbos {
			log.Debug("filelist: sub_id: %d", dbo.RefSubFile)
			if dbo.RefSubFile > 0 {
				subf, err := SubFileGet(f.RefSubFile)
				if err == nil {
					jf.SubFile = subf
				}
			}
			ret = append(ret, dbo)
		}*/
	return dbos, nil
}

func FileGet(id int) (*DbFile, error) {
	var o DbFile
	err := db.One("ID", id, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func FileGetByUrl(url string) (*DbFile, error) {
	var o DbFile
	err := db.One("UrlPath", url, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func FileGetByRedirectUrl(url string) (*DbFile, error) {
	var o DbFile
	err := db.One("RedirectPath", url, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func FileDirExists(url string) bool {
	var o []DbFile
	if url == "" {
		return false
	}
	if url[len(url)-1] != '/' {
		url += "/"
	}
	err := db.Prefix("UrlPath", url, &o)
	if err != nil {
		return false
	}
	return true
}

func FileDelete(id int) error {
	f := &DbFile{
		ID: id,
	}
	err := db.DeleteStruct(f)
	if err != nil {
		return err
	}
	return nil
}

func FileUpdate(id int, o *DbFile) (*DbFile, error) {
	if err := db.Update(&DbFile{ID: id, Name: o.Name, UrlPath: o.UrlPath, MimeType: o.MimeType, RefSubFile: o.RefSubFile, SubName: o.SubName, RedirectPath: o.RedirectPath, SubMimeType: o.SubMimeType}); err != nil {
		return nil, err
	}
	if err := db.UpdateField(&DbFile{ID: id}, "RedirectPath", o.RedirectPath); err != nil {
		return nil, err
	}
	return o, nil
}

func FileResetSubFile(id int) (*DbFile, error) {
	if err := db.UpdateField(&DbFile{ID: id}, "RefSubFile", 0); err != nil {
		return nil, err
	}
	o, err := FileGet(id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func FileEnable(id int, enable bool) (*DbFile, error) {
	if err := db.UpdateField(&DbFile{ID: id}, "IsEnabled", enable); err != nil {
		return nil, err
	}
	o, err := FileGet(id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func FilePause(id int, pause bool) (*DbFile, error) {
	if err := db.UpdateField(&DbFile{ID: id}, "IsPaused", pause); err != nil {
		return nil, err
	}
	o, err := FileGet(id)
	if err != nil {
		return nil, err
	}
	return o, nil
}
