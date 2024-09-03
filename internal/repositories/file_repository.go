package repositories

import (
	context2 "context"
	"database/sql"
	"famdocs-zoo/models"
)

type fileRepository struct {
	db *sql.DB
}

type FileRepository interface {
	Create(file *models.File) (*models.File, error)
	Delete(file *models.File) (*models.File, error)
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &fileRepository{db}
}

func (fr *fileRepository) Create(file *models.File) (*models.File, error) {
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryResponse, err :=
		txn.Query(`insert into files(name,folder_id,path) values($1,$2,$3) returning id`,
			file.Name, file.FolderId, file.Path)
	if err != nil {
		return nil, err
	}
	var fileId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&fileId)
		if err != nil {
			return nil, err
		}
	}
	file.Id = fileId
	queryResponse.Close()
	queryResponse, err =
		txn.Query(`insert into folder_to_file_map(parent_folder_id,file_id) values($1,$2)`,
			file.FolderId, fileId)
	if err != nil {
		return nil, err
	}
	queryResponse.Close()
	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	return file, err
}

func (fr *fileRepository) Delete(file *models.File) (*models.File, error) {
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	_, err = txn.Query(
		`delete from folder_to_file_map where file_id=?`, file.Id)
	if err != nil {
		return nil, err
	}
	_, err = txn.Query(
		`delete from files where id=?`, file.Id)
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	return file, nil
}
