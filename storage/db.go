package storage

import (
	"os"
	"path/filepath"

	"github.com/asdine/storm"

	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/utils"
)

var db *storm.DB

func Open(path string) error {
	var err error

	err = os.MkdirAll(filepath.Dir(path), 0600)
	if err != nil {
		return err
	}

	db, err = storm.Open(path)
	if err != nil {
		return err
	}

	err = db.Init(&DbFile{})
	if err != nil {
		return err
	}
	err = db.Init(&DbSubFile{})
	if err != nil {
		return err
	}
	err = db.Init(&DbUser{})
	if err != nil {
		return err
	}
	err = db.Init(&DbSession{})
	if err != nil {
		return err
	}
	err = db.Init(&DbConfig{})
	if err != nil {
		return err
	}

	// initialize config
	err = initConfig()
	if err != nil {
		return err
	}

	return nil
}

func initConfig() error {
	o, err := ConfigGet(1)
	if err != nil {
		o = &DbConfig{
			ID:          1,
			SecretPath:  "/pwndrop",
			RedirectUrl: "https://www.youtube.com/watch?v=oHg5SJYRHA0",
			CookieName:  utils.GenRandomString(4),
			CookieToken: utils.GenRandomHash(),
		}
		_, err = ConfigCreate(o)
		if err != nil {
			return err
		}
	}
	log.Debug("secret_path: %s", o.SecretPath)
	// update added values here in future updates
	return nil
}
