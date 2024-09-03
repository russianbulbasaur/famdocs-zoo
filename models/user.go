package models

type User struct {
	Name   string `json:"name"`
	Id     int64  `json:"user_id"`
	Avatar int8   `json:"avatar"`
	Phone  string `json:"phone"`
	Token  string `json:"token"`
}
