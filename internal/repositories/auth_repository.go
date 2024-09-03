package repositories

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"famdocs-zoo/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type authRepository struct {
	db *sql.DB
}

type AuthRepository interface {
	Login(string, string, string) (*models.User, error)
	Signup(user *models.User) (*models.User, error)
}

type FirebaseResponse struct {
	Phone string `json:"phoneNumber"`
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func userInDatabase(phone string, db *sql.DB) (*models.User, error) {
	var user models.User
	result, err := db.Query(`select id,name,phone,avatar from users where phone=$1`, phone)
	defer result.Close()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	if !result.Next() {
		return nil, err
	}
	err = result.Scan(&user.Id, &user.Name, &user.Phone, &user.Avatar)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &user, nil
}

func (ar *authRepository) Login(phone, token, otp string) (*models.User, error) {
	user, err := userInDatabase(phone, ar.db)
	if !firebaseAuth(token, otp, fmt.Sprintf("+91%s", phone)) {
		return nil, errors.New("invalid firebase token")
	}
	return user, err
}

func (ar *authRepository) Signup(user *models.User) (*models.User, error) {
	dbUser, err := userInDatabase(user.Phone, ar.db)
	if dbUser != nil {
		return nil, errors.New("user already in db")
	}
	response, err := ar.db.Query(`insert into users(name,phone,avatar) values($1,$2,$3) returning id`,
		user.Name, user.Phone, user.Avatar)
	defer response.Close()
	if err != nil {
		return nil, err
	}
	var userId int64
	if response.Next() {
		err = response.Scan(&userId)
		if err != nil {
			return nil, err
		}
	}
	user.Id = userId
	return user, nil
}

func firebaseAuth(firebaseToken, otp, phone string) bool {
	apiKey := os.Getenv("FIREBASE_API_KEY")
	body, _ := json.Marshal(map[string]string{
		"sessionInfo": firebaseToken,
		"code":        otp,
	})
	urlString := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPhoneNumber?key=%s", apiKey)
	client := &http.Client{}
	response, err := client.Post(urlString, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return false
	}
	defer response.Body.Close()
	body, _ = io.ReadAll(response.Body)
	var firebaseRes FirebaseResponse
	err = json.Unmarshal(body, &firebaseRes)
	if err != nil {
		log.Println(err)
		return false
	}
	if firebaseRes.Phone == phone {
		return true
	}
	return false
}
