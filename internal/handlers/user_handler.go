package handlers

import (
	"encoding/json"
	"famdocs-zoo/helpers"
	"famdocs-zoo/internal/services"
	"famdocs-zoo/models"
	"net/http"
)

type userHandler struct {
	userServ services.UserService
}

type UserHandler interface {
	Update(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetUserFamilies(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(userServ services.UserService) UserHandler {
	return &userHandler{userServ}
}

func (uh *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 0, err)
		return
	}
	var user *models.User
	err = json.Unmarshal([]byte(r.Form.Get("user")), &user)
	if err != nil || user == nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	response, err := uh.userServ.GetUserFromId(user.Id)
	if err != nil {
		helpers.InternalError(w, 2, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 3, err)
	}
}

func (uh *userHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (uh *userHandler) GetUserFamilies(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	var user models.User
	userString := r.Form.Get("user")
	err = json.Unmarshal([]byte(userString), &user)
	if err != nil {
		helpers.InvalidParametersError(w, -1, err)
		return
	}
	response, err := uh.userServ.GetUserFamilies(&user)
	if err != nil {
		helpers.InternalError(w, 2, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 3, err)
	}
}
