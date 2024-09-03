package handlers

import (
	"encoding/json"
	"errors"
	"famdocs-zoo/helpers"
	"famdocs-zoo/internal/services"
	"famdocs-zoo/models"
	"net/http"
	"strconv"
)

type familyHandler struct {
	familyService services.FamilyService
}

type FamilyHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetRootFolder(w http.ResponseWriter, r *http.Request)
	JoinFamily(w http.ResponseWriter, r *http.Request)
}

func NewFamilyHandler(familyService services.FamilyService) FamilyHandler {
	return &familyHandler{familyService}
}

type CreateRequest struct {
	UserId     int64  `json:"user_id"`
	FamilyName string `json:"family_name"`
	ShaHash    string `json:"sha_hash"`
}

func (fh *familyHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	var request CreateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.FormParseError(w, 2, err)
		return
	}
	var user models.User
	err = json.Unmarshal([]byte(r.Form.Get("user")), &user)
	if err != nil {
		helpers.FormParseError(w, 3, err)
		return
	}
	res, err := fh.familyService.Create(request.UserId, request.FamilyName,
		request.ShaHash, user.Name)
	if err != nil {
		helpers.InternalError(w, 5, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		helpers.WriteError(w, 6, err)
		return
	}
}

func (fh *familyHandler) GetRootFolder(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	familyId, err := strconv.ParseInt(params.Get("family_id"), 10, 64)
	if err != nil {
		helpers.InvalidParametersError(w, 1, err)
		return
	}
	response, err := fh.familyService.GetRootFolder(familyId)
	if err != nil {
		helpers.InternalError(w, 2, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 3, err)
		return
	}
}

func (fh *familyHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	memberId, err := strconv.ParseInt(r.Form.Get("member_id"), 10, 64)
	if err != nil {
		helpers.InvalidParametersError(w, 2, err)
		return
	}
	familyId, err := strconv.ParseInt(r.Form.Get("family_id"), 10, 64)
	if err != nil {
		helpers.InvalidParametersError(w, 3, err)
		return
	}
	var user *models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user == nil {
		helpers.InvalidParametersError(w, 4, err)
		return
	}
	response, err := fh.familyService.AddMember(user, memberId, familyId)
	if err != nil {
		helpers.InternalError(w, 5, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 6, err)
		return
	}
}

func (fh *familyHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	familyId, err := strconv.ParseInt(r.Form.Get("family_id"), 10, 64)
	if err != nil {
		helpers.InvalidParametersError(w, 2, err)
		return
	}
	response, err := fh.familyService.GetMembers(familyId)
	if err != nil {
		helpers.InternalError(w, 3, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 4, err)
	}
}

type JoinRequest struct {
	JoinCode string `json:"join_code"`
	Password string `json:"password"`
}

func (fh *familyHandler) JoinFamily(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	var user models.User
	err = json.Unmarshal([]byte(r.Form.Get("user")), &user)
	if err != nil {
		helpers.InvalidParametersError(w, 2, err)
		return
	}
	var request JoinRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	joinCode := request.JoinCode
	password := request.Password
	if joinCode == "" || password == "" {
		helpers.MissingParametersError(w, 2, errors.New("invalid params"))
		return
	}
	familyId, err := fh.familyService.GetFamilyIDFromJoinCode(joinCode)
	if err != nil {
		helpers.InvalidParametersError(w, 3, err)
		return
	}
	response, err := fh.familyService.JoinFamily(familyId, &user, password)
	if err != nil {
		helpers.InternalError(w, 4, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 5, err)
		return
	}
}
