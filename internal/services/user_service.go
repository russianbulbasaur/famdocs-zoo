package services

import (
	"encoding/json"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/models"
	"log"
)

type userService struct {
	userRepo repositories.UserRepository
}

type UserService interface {
	GetUserFromId(int64) ([]byte, error)
	GetUserFamilies(*models.User) ([]byte, error)
}

func NewUserService(userServ repositories.UserRepository) UserService {
	return &userService{userServ}
}

func (us *userService) GetUserFromId(userId int64) ([]byte, error) {
	user, err := us.userRepo.GetUserFromId(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, err
}

func (us *userService) GetUserFamilies(user *models.User) ([]byte, error) {
	familyList, err := us.userRepo.GetUserFamilies(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response, err := json.Marshal(familyList)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, err
}
