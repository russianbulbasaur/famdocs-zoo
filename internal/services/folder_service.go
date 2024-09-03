package services

import (
	"encoding/json"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/models"
	"log"
)

type folderService struct {
	folderRepo repositories.FolderRepository
}

type FolderService interface {
	Create(folder *models.Folder) ([]byte, error)
	GetContents(folderId int64) ([]byte, error)
}

func NewFolderService(fr repositories.FolderRepository) FolderService {
	return &folderService{fr}
}

func (fs *folderService) Create(folder *models.Folder) ([]byte, error) {
	folder, err := fs.folderRepo.Create(folder)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := json.Marshal(folder)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (fs *folderService) GetContents(folderId int64) ([]byte, error) {
	content, err := fs.folderRepo.GetContents(folderId)
	print()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response, err := json.Marshal(content)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, err
}
