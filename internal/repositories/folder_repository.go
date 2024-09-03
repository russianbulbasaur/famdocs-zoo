package repositories

import (
	context2 "context"
	"database/sql"
	"famdocs-zoo/models"
)

type folderRepository struct {
	db *sql.DB
}

type FolderRepository interface {
	Create(folder *models.Folder) (*models.Folder, error)
	GetContents(folderId int64) (*FolderContent, error)
}

func NewFolderRepository(db *sql.DB) FolderRepository {
	return &folderRepository{db}
}

func (fr *folderRepository) Create(folder *models.Folder) (*models.Folder, error) {
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryResponse, err :=
		txn.Query(`insert into folders(name,family_id) values($1,$2) returning id`,
			folder.Name, folder.FamilyId)
	if err != nil {
		return nil, err
	}
	var folderId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&folderId)
		if err != nil {
			return nil, err
		}
	}
	folder.Id = folderId
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}
	queryResponse, err =
		txn.Query(`insert into folder_to_folder_map(folder_id,parent_folder_id) values($1,$2)`,
			folderId, folder.ParentFolderId)
	if err != nil {
		return nil, err
	}
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	return folder, err
}

type FolderContent struct {
	Files   []*models.File   `json:"files"`
	Folders []*models.Folder `json:"folders"`
}

func (fr *folderRepository) GetContents(folderId int64) (*FolderContent, error) {
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryResponse, err :=
		txn.Query(
			`select folders.id as id ,
                   folders.name as name,folders.family_id as family_id,
                   ffm.parent_folder_id as parent_folder_id from folders inner join 
                   folder_to_folder_map ffm on ffm.folder_id=folders.id where ffm.parent_folder_id=$1`,
			folderId)
	if err != nil {
		return nil, err
	}
	var subFolderList []*models.Folder
	for queryResponse.Next() {
		subFolder := new(models.Folder)
		err = queryResponse.Scan(&subFolder.Id,
			&subFolder.Name,
			&subFolder.FamilyId, &subFolder.ParentFolderId)
		if err != nil {
			return nil, err
		}
		subFolderList = append(subFolderList, subFolder)
	}
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}

	queryResponse, err =
		txn.Query(
			`select files.id as id ,
                   files.name as name,files.path as path,
                   files.folder_id as folder_id from files inner join 
                   folder_to_file_map ffm on ffm.file_id=files.id where ffm.parent_folder_id=$1`,
			folderId)
	if err != nil {
		return nil, err
	}
	var subFilesList []*models.File
	for queryResponse.Next() {
		subFile := new(models.File)
		err = queryResponse.Scan(&subFile.Id, &subFile.Name, &subFile.Path, &subFile.FolderId)
		if err != nil {
			return nil, err
		}
		subFilesList = append(subFilesList, subFile)
	}
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}
	content := new(FolderContent)
	content.Files = subFilesList
	content.Folders = subFolderList
	return content, err
}
