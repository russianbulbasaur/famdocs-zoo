package models

type Member struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
	Avatar int8   `json:"avatar"`
}
