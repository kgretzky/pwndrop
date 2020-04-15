package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	//"github.com/gorilla/mux"

	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
	"github.com/kgretzky/pwndrop/utils"
)

const AUTH_COOKIE_NAME = "t"
const AUTH_SESSION_TIMEOUT_SECS = 24 * 60 * 60

func AuthOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
}

func AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	type AuthResponse struct {
		Status int `json:"status"`
	}

	users, err := storage.UserList()
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	resp := &AuthResponse{}
	if len(users) == 0 {
		resp.Status = 0
		DumpResponse(w, "ok", http.StatusOK, 0, resp)
		return
	}

	_, err = AuthSession(r)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	resp.Status = 1
	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type LoginResponse struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}

	j := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	log.Debug("username: %s", j.Username)

	o, err := storage.UserGetByName(j.Username)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(o.Password), []byte(j.Password))
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	token := utils.GenRandomHash()
	s := &storage.DbSession{
		Uid:        o.ID,
		Token:      token,
		CreateTime: time.Now().Unix(),
	}

	_, err = storage.SessionCreate(s)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	resp := &LoginResponse{
		Username: o.Name,
		Token:    token,
	}

	ck := &http.Cookie{
		Domain:   "",
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Name:     AUTH_COOKIE_NAME,
		Value:    token,
	}
	http.SetCookie(w, ck)

	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	ck, err := r.Cookie(AUTH_COOKIE_NAME)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	token := ck.Value

	s, err := storage.SessionGetByToken(token)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	err = storage.SessionDelete(s.ID)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	deleteCookie(AUTH_COOKIE_NAME, w)
	DumpResponse(w, "ok", http.StatusOK, 0, nil)
}

func ClearSecretSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie_name := Cfg.GetCookieName()
	deleteCookie(cookie_name, w)
	DumpResponse(w, "ok", http.StatusOK, 0, nil)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type CreateUserResponse struct {
		Username string `json:"username"`
	}

	users, err := storage.UserList()
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusInternalServerError, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	_, err = AuthSession(r)
	if len(users) > 0 && err != nil {
		DumpResponse(w, err.Error(), http.StatusUnauthorized, API_ERROR_BAD_AUTHENTICATION, nil)
		return
	}

	j := CreateUserRequest{}
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	if j.Username == "" || j.Password == "" {
		DumpResponse(w, "bad request", http.StatusBadRequest, API_ERROR_BAD_REQUEST, nil)
		return
	}

	_, err = storage.UserGetByName(j.Username)
	if err == nil {
		DumpResponse(w, "user already exists", http.StatusOK, API_ERROR_USER_ALREADY_EXISTS, nil)
		return
	}

	phash, err := bcrypt.GenerateFromPassword([]byte(j.Password), 10)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	o := &storage.DbUser{
		Name:     j.Username,
		Password: string(phash),
	}

	_, err = storage.UserCreate(o)
	if err != nil {
		DumpResponse(w, err.Error(), http.StatusBadRequest, API_ERROR_FILE_DATABASE_FAILED, nil)
		return
	}

	resp := &CreateUserResponse{
		Username: j.Username,
	}
	DumpResponse(w, "ok", http.StatusOK, 0, resp)
}

func AuthSession(r *http.Request) (int, error) {
	ck, err := r.Cookie(AUTH_COOKIE_NAME)
	if err != nil {
		return -1, err
	}

	token := ck.Value

	s, err := storage.SessionGetByToken(token)
	if err != nil {
		return -1, err
	}

	if time.Now().After(time.Unix(s.CreateTime, 0).Add(AUTH_SESSION_TIMEOUT_SECS * time.Second)) {
		storage.SessionDelete(s.ID)
		return -1, fmt.Errorf("session token expired")
	}

	return s.Uid, nil
}

func deleteCookie(name string, w http.ResponseWriter) {
	ck := &http.Cookie{
		Domain:   "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Name:     name,
		Value:    "",
	}
	http.SetCookie(w, ck)
}
