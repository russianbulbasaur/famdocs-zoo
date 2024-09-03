package services

import (
	"encoding/json"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/models"
	"log"
	"strconv"
)

type familyService struct {
	famRepo repositories.FamilyRepository
}

type FamilyService interface {
	Create(int64, string, string, string) ([]byte, error)
	GetRootFolder(int64) ([]byte, error)
	AddMember(*models.User, int64, int64) ([]byte, error)
	GetMembers(int64) ([]byte, error)
	JoinFamily(int64, *models.User, string) ([]byte, error)
	GetFamilyIDFromJoinCode(string) (int64, error)
}

func NewFamilyService(famRepo repositories.FamilyRepository) FamilyService {
	return &familyService{famRepo}
}

func (fs *familyService) Create(userId int64, familyName string, shaHash string,
	userName string) ([]byte, error) {
	family, err := fs.famRepo.Create(userId, familyName, shaHash, userName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	familyResponse, err := json.Marshal(family)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return familyResponse, err
}

func (fs *familyService) GetRootFolder(familyId int64) ([]byte, error) {
	folder, err := fs.famRepo.GetRootFolder(familyId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response, err := json.Marshal(folder)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, err
}

func (fs *familyService) AddMember(user *models.User, memberId int64, familyId int64) ([]byte, error) {
	addResponse, err := fs.famRepo.AddMember(user, memberId, familyId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if addResponse {
		return json.Marshal(map[string]string{
			"message": "Added member",
		})
	}
	return nil, err
}

func (fs *familyService) GetMembers(familyId int64) ([]byte, error) {
	memberList, err := fs.famRepo.GetMembers(familyId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response, err := json.Marshal(memberList)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (fs *familyService) JoinFamily(familyId int64, user *models.User,
	password string) ([]byte, error) {
	response, err := fs.famRepo.JoinFamily(familyId, user, password)
	if err != nil {
		return nil, err
	}
	return []byte(strconv.FormatBool(response)), err
}

func (fs *familyService) GetFamilyIDFromJoinCode(joinCode string) (int64, error) {
	id := ""
	for i := 3; i < len(joinCode); i++ {
		id += string(joinCode[i])
	}
	return strconv.ParseInt(id, 10, 64)
}
