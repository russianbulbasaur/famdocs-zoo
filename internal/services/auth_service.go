package services

import (
	"encoding/json"
	"errors"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/models"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"
)

type authService struct {
	authRepo repositories.AuthRepository
}

type AuthService interface {
	Login(string, string, string) ([]byte, error)
	Signup(user *models.User) ([]byte, error)
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{authRepo}
}

func (as *authService) Login(phone, token, otp string) ([]byte, error) {
	user, err := as.authRepo.Login(phone, token, otp)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.Token = makeToken(user)
		return json.Marshal(user)
	} else {
		var newUser models.User
		newUser.Phone = phone
		newUser.Token = makeSignupToken(phone)
		return json.Marshal(newUser)
	}
}

func (as *authService) Signup(user *models.User) ([]byte, error) {
	if !verifySignupToken(user.Token, user.Phone) {
		log.Println("invalid token")
		return nil, errors.New("invalid token")
	}
	user, err := as.authRepo.Signup(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	user.Token = makeToken(user)
	return json.Marshal(user)
}

func verifySignupToken(token, phone string) bool {
	data, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		key := []byte(os.Getenv("JWT_SECRET"))
		return key, nil
	})
	if err != nil || !data.Valid {
		log.Println(err)
		return false
	}
	sub, err := data.Claims.GetSubject()
	if err != nil || sub != phone {
		log.Println(err)
		return false
	}
	return true
}

func makeSignupToken(phone string) string {
	key := []byte(os.Getenv("JWT_SECRET"))
	data := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": phone,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	token, err := data.SignedString(key)
	if err != nil {
		log.Println(err)
	}
	return token
}

func makeToken(user *models.User) string {
	key := []byte(os.Getenv("JWT_SECRET"))
	data := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(user.Id, 10),
	})
	token, err := data.SignedString(key)
	if err != nil {
		log.Println(err)
	}
	return token
}
