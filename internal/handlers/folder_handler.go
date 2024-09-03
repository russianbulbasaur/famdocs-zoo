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

type folderHandler struct {
	folderService services.FolderService
}

type FolderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetContents(w http.ResponseWriter, r *http.Request)
}

func NewFolderHandler(fs services.FolderService) FolderHandler {
	return &folderHandler{fs}
}

func (fh *folderHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	var folder models.Folder
	err = json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		helpers.InvalidParametersError(w, 2, err)
		return
	}
	response, err := fh.folderService.Create(&folder)
	if err != nil {
		helpers.InternalError(w, 3, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		helpers.WriteError(w, 4, err)
		return
	}
}

func (fh *folderHandler) GetContents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	folderId, err := strconv.ParseInt(params.Get("folder_id"), 10, 64)
	if err != nil {
		helpers.InvalidParametersError(w, 1, errors.New("empty paramter"))
		return
	}
	response, err := fh.folderService.GetContents(folderId)
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
