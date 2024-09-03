package models

type File struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	FolderId int64  `json:"folder_id"`
	Path     string `json:"path"`
}
