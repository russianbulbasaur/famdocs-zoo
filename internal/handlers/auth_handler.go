package handlers

import (
	"encoding/json"
	"errors"
	"famdocs-zoo/helpers"
	"famdocs-zoo/internal/services"
	"famdocs-zoo/models"
	"net/http"
)

type authHandler struct {
	authServ services.AuthService
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(authServ services.AuthService) AuthHandler {
	return &authHandler{authServ}
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	phone := params.Get("phone")
	otp := params.Get("otp")
	token := params.Get("token")
	if phone == "" || otp == "" || token == "" {
		helpers.MissingParametersError(w, 1, errors.New("missing params"))
		return
	}
	res, err := ah.authServ.Login(phone, token, otp)
	if err != nil {
		helpers.InternalError(w, 2, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		helpers.WriteError(w, 3, err)
	}
}

func (ah *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		helpers.MissingParametersError(w, 1, err)
		return
	}
	res, err := ah.authServ.Signup(user)
	if err != nil {
		helpers.InternalError(w, 2, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		helpers.WriteError(w, 3, err)
	}
}
