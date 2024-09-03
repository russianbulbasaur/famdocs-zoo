package models

type Folder struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FamilyId       int64  `json:"family_id"`
	ParentFolderId int64  `json:"parent_folder_id"`
}
