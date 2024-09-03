package handlers

import (
	"encoding/json"
	"errors"
	"famdocs-zoo/helpers"
	"famdocs-zoo/internal/services"
	"famdocs-zoo/models"
	uuid2 "github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type fileHandler struct {
	fileService services.FileService
}

type FileHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Download(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

func NewFileHandler(fs services.FileService) FileHandler {
	return &fileHandler{fs}
}

func (fh *fileHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(800)
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	uploadedFile, _, err := r.FormFile("file")
	if uploadedFile == nil {
		helpers.FormParseError(w, 1, errors.New("file is empty"))
		return
	}
	uuid := uuid2.New()
	filePointer, err := os.Create(filepath.Join("uploads/", uuid.String()))
	if filePointer == nil {
		helpers.FormParseError(w, 1, errors.New("file pointer error"))
		return
	}
	_, err = io.Copy(filePointer, uploadedFile)
	if err != nil {
		helpers.FormParseError(w, 1, errors.New("writing error"))
		return
	}
	fileModel := new(models.File)
	fileModel.FolderId, err = strconv.ParseInt(r.FormValue("folder_id"), 10, 64)
	fileModel.Path = uuid.String()
	fileModel.Name = r.FormValue("file_name") + "." + r.FormValue("ext")
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	response, err := fh.fileService.Create(fileModel)
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

func (fh *fileHandler) Download(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filePath := params.Get("path")
	http.ServeFile(w, r, "uploads/"+filePath)
}

func (fh *fileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var file *models.File
	err := json.NewDecoder(r.Body).Decode(file)
	if err != nil {
		helpers.FormParseError(w, 1, err)
		return
	}
	response, err := fh.fileService.Delete(file)
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
