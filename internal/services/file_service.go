package services

import (
	"encoding/json"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/models"
	"log"
)

type fileService struct {
	fileRepo repositories.FileRepository
}

type FileService interface {
	Create(file *models.File) ([]byte, error)
	Delete(file *models.File) ([]byte, error)
}

func NewFileService(fr repositories.FileRepository) FileService {
	return &fileService{fr}
}

func (fs *fileService) Create(file *models.File) ([]byte, error) {
	file, err := fs.fileRepo.Create(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := json.Marshal(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, err
}

func (fs *fileService) Delete(file *models.File) ([]byte, error) {
	response, err := fs.fileRepo.Delete(file)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return res, err
}
