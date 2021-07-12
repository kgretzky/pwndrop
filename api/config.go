package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"net"
	"fmt"

	"github.com/kgretzky/pwndrop/storage"
	"github.com/kgretzky/pwndrop/utils"
)

func ConfigOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
}

func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	o, err := storage.ConfigGet(1)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, o)
}

func ConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// #### CHECK IF AUTHENTICATED ####
	_, err := AuthSession(r)
	if err != nil {
		DumpResponse(w, "unauthorized", http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	old_cfg, err := storage.ConfigGet(1)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	o := storage.DbConfig{}
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	if o.SecretPath == "" || o.CookieName == "" || o.CookieToken == "" {
		DumpResponse(w, "missing config variables", http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	ips := ""
	if o.UsrWhiteList != "" && o.UsrBlackList != "" {
		ips = o.UsrWhiteList + "," + o.UsrBlackList
	} else if o.UsrWhiteList != "" {
		ips = o.UsrWhiteList
	} else if o.UsrBlackList != "" {
		ips = o.UsrBlackList
	}
	for _, ip := range strings.Split(ips, ",") {
		if ip == "" {
			continue
		}
		_, _, err1 := net.ParseCIDR(ip)
		err2 := net.ParseIP(ip)
		if err1 != nil && err2 == nil {
			DumpResponse(w, fmt.Sprintf("invalid IP: %s", ip) , http.StatusBadRequest, API_ERROR_INVALID_IP, nil)
			return
		}
	}

	if o.SecretPath[0] != '/' {
		o.SecretPath = "/" + o.SecretPath
	}
	if o.SecretPath != old_cfg.SecretPath {
		o.CookieName = utils.GenRandomString(4)
		o.CookieToken = utils.GenRandomHash()
	}

	ret, err := storage.ConfigUpdate(1, &o)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}
	DumpResponse(w, "ok", http.StatusOK, 0, ret)
}
