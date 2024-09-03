package models

type Family struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Owner   int64  `json:"owner"`
	ShaHash string `json:"sha_hash"`
	Root    Folder `json:"root"`
}
