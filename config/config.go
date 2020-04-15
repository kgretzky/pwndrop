package config

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ini.v1"
	"path/filepath"
	"strconv"

	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
	"github.com/kgretzky/pwndrop/utils"
)

const (
	INI_SERVER         = "pwndrop"
	INI_VAR_LISTEN_IP  = "listen_ip"
	INI_VAR_HTTP_PORT  = "http_port"
	INI_VAR_HTTPS_PORT = "https_port"
	INI_VAR_DATA_DIR   = "data_dir"
	INI_VAR_ADMIN_DIR  = "admin_dir"

	INI_SETUP              = "setup"
	INI_SETUP_USERNAME     = "username"
	INI_SETUP_PASSWORD     = "password"
	INI_SETUP_REDIRECT_URL = "redirect_url"
	INI_SETUP_SECRET_PATH  = "secret_path"
)

type Config struct {
	ini      *ini.File
	path     string
	exec_dir string
}

func NewConfig(path string) (*Config, error) {
	var err error
	c := &Config{
		path:     path,
		exec_dir: utils.GetExecDir(),
	}

	data_dir := filepath.Join(c.exec_dir, "data")
	admin_dir := filepath.Join(c.exec_dir, "admin")

	defs := map[string]string{
		INI_VAR_LISTEN_IP:  "",
		INI_VAR_HTTP_PORT:  "80",
		INI_VAR_HTTPS_PORT: "443",
		INI_VAR_DATA_DIR:   data_dir,
		INI_VAR_ADMIN_DIR:  admin_dir,
	}

	c.ini, err = ini.Load(path)
	if err != nil {
		log.Warning("config file not found at path: %s", path)
		c.ini = ini.Empty()
	}

	if _, err = c.ini.GetSection(INI_SERVER); err != nil {
		c.ini.NewSection(INI_SERVER)
	}

	for k, v := range defs {
		if _, err = c.ini.Section(INI_SERVER).GetKey(k); err != nil {
			c.ini.Section(INI_SERVER).NewKey(k, v)
		}
	}

	return c, nil
}

func (c *Config) HandleSetup() error {
	if _, err := c.ini.GetSection(INI_SETUP); err == nil {
		o, err := storage.ConfigGet(1)
		if err != nil {
			log.Error("config: can't get config from db")
		}

		var username, password, redirect_url, secret_path string

		if k, err := c.ini.Section(INI_SETUP).GetKey(INI_SETUP_USERNAME); err == nil {
			username = k.String()
		}
		if k, err := c.ini.Section(INI_SETUP).GetKey(INI_SETUP_PASSWORD); err == nil {
			password = k.String()
		}
		if k, err := c.ini.Section(INI_SETUP).GetKey(INI_SETUP_REDIRECT_URL); err == nil {
			redirect_url = k.String()
			o.RedirectUrl = redirect_url
			log.Important("setup: redirect url set to: %s", redirect_url)
		}
		if k, err := c.ini.Section(INI_SETUP).GetKey(INI_SETUP_SECRET_PATH); err == nil {
			secret_path = k.String()
			if secret_path[0] != '/' {
				secret_path = "/" + secret_path
			}
			if len(secret_path) >= 2 {
				o.CookieName = utils.GenRandomString(4)
				o.CookieToken = utils.GenRandomHash()
				o.SecretPath = secret_path
				log.Important("setup: secret path set to: %s", secret_path)
			}
		}

		if len(username) > 0 && len(password) > 0 {
			phash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
			if err == nil {
				o := &storage.DbUser{
					Name:     username,
					Password: string(phash),
				}

				storage.UserDelete(1)
				_, err = storage.UserCreate(o)
				if err == nil {
					log.Important("setup: created user account: %s", username)
				} else {
					log.Error("setup: failed to create user account: %s", err)
				}
				err = storage.SessionDeleteAll()
				if err != nil {
					log.Error("failed to delete active sessions: %s", err)
				}
			}
		}

		_, err = storage.ConfigUpdate(1, o)
		if err != nil {
			log.Error("config: can't save config to db")
		}

		c.ini.DeleteSection(INI_SETUP)
	}
	return nil
}

func (c *Config) Save() error {
	err := c.ini.SaveTo(c.path)
	if err != nil {
		return fmt.Errorf("failed to save config file")
	}
	return nil
}

func (c *Config) GetListenIP() string {
	s, _ := c.Get(INI_VAR_LISTEN_IP)
	return s
}

func (c *Config) GetHttpPort() int {
	s, _ := c.Get(INI_VAR_HTTP_PORT)
	port, _ := strconv.Atoi(s)
	return port
}

func (c *Config) GetHttpsPort() int {
	s, _ := c.Get(INI_VAR_HTTPS_PORT)
	port, _ := strconv.Atoi(s)
	return port
}

func (c *Config) GetSecretPath() string {
	o, err := storage.ConfigGet(1)
	if err != nil {
		return ""
	}
	return o.SecretPath
}

func (c *Config) GetDataDir() string {
	dir, _ := c.Get(INI_VAR_DATA_DIR)
	return c.joinPath(c.exec_dir, dir)
}

func (c *Config) GetAdminDir() string {
	dir, _ := c.Get(INI_VAR_ADMIN_DIR)
	return c.joinPath(c.exec_dir, dir)
}

func (c *Config) GetCookieName() string {
	o, err := storage.ConfigGet(1)
	if err != nil {
		return ""
	}
	return o.CookieName
}

func (c *Config) GetCookieToken() string {
	o, err := storage.ConfigGet(1)
	if err != nil {
		return ""
	}
	return o.CookieToken
}

func (c *Config) GetRedirectUrl() string {
	o, err := storage.ConfigGet(1)
	if err != nil {
		return ""
	}
	return o.RedirectUrl
}

func (c *Config) Get(key string) (string, error) {
	section, err := c.ini.GetSection(INI_SERVER)
	if err != nil {
		return "", err
	}
	if section.HasKey(key) {
		return section.Key(key).String(), nil
	}
	return "", fmt.Errorf("config key '%s' not found", key)
}

func (c *Config) Set(key string, value string) error {
	section, err := c.ini.GetSection(INI_SERVER)
	if err != nil {
		return err
	}
	if section.HasKey(key) {
		section.Key(key).SetValue(value)
	} else {
		section.NewKey(key, value)
	}
	err = c.ini.SaveTo(c.path)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) joinPath(base_path string, rel_path string) string {
	var ret string
	if filepath.IsAbs(rel_path) {
		ret = rel_path
	} else {
		ret = filepath.Join(base_path, rel_path)
	}
	return ret
}
