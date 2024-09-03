package repositories

import (
	context2 "context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"famdocs-zoo/models"
)

type familyRepository struct {
	db *sql.DB
}

type FamilyRepository interface {
	Create(int64, string, string, string) (*models.Family, error)
	GetMembers(int64) ([]*models.Member, error)
	AddMember(*models.User, int64, int64) (bool, error)
	GetRootFolder(int64) (*models.Folder, error)
	JoinFamily(int64, *models.User, string) (bool, error)
}

func NewFamilyRepository(db *sql.DB) FamilyRepository {
	return &familyRepository{db}
}

func (fr *familyRepository) Create(userId int64, familyName string, shaHash string,
	userName string) (*models.Family, error) {
	family := new(models.Family)
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryResponse, err :=
		txn.Query(`insert into families(name,owner,sha_hash) values($1,$2,$3) returning id;`,
			familyName, userId, shaHash)
	if err != nil {
		return nil, err
	}
	var familyId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&familyId)
		if err != nil {
			return nil, err
		}
	}
	family.Owner = userId
	family.Name = familyName
	family.Id = familyId
	family.ShaHash = shaHash
	queryResponse.Close()
	queryResponse, err = txn.Query(
		`insert into user_family_map(user_id,family_id) values($1,$2)`,
		userId, familyId)
	if err != nil {
		return nil, err
	}
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}
	queryResponse, err = txn.Query(`insert into folders(name,family_id) values($1,$2) returning id`,
		"root", familyId)
	if err != nil {
		return nil, err
	}
	rootFolder := new(models.Folder)
	var rootFolderId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&rootFolderId)
		if err != nil {
			return nil, err
		}
	}
	rootFolder.FamilyId = familyId
	rootFolder.Name = "root"
	rootFolder.Id = rootFolderId
	family.Root = *rootFolder
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}

	queryResponse, err =
		txn.Query(`insert into folders(name,family_id) values($1,$2) returning id`,
			userName+"'s Folder", family.Id)
	if err != nil {
		return nil, err
	}
	var ownerFolderId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&ownerFolderId)
		if err != nil {
			return nil, err
		}
	}
	err = queryResponse.Close()
	if err != nil {
		return nil, err
	}
	queryResponse, err =
		txn.Query(`insert into folder_to_folder_map(folder_id,parent_folder_id) values($1,$2)`,
			ownerFolderId, rootFolderId)
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
	return family, nil
}

func (fr *familyRepository) AddMember(user *models.User, memberId int64, familyId int64) (bool, error) {
	dbResponse, err := fr.db.Query(`select owner from families where id=$1`, familyId)
	defer dbResponse.Close()
	if err != nil {
		return false, err
	}
	var owner int64
	if dbResponse.Next() {
		err = dbResponse.Scan(owner)
		if err != nil {
			return false, err
		}
	}
	if owner != user.Id {
		return false, errors.New("User does not own the family")
	}
	dbResponse, err = fr.db.Query(`insert into user_family_map(user_id,family_id) values($1,$2)`,
		memberId, familyId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fr *familyRepository) GetMembers(familyId int64) ([]*models.Member, error) {
	dbResponse, err := fr.db.Query(`select * from user_family_map where family_id=$1`, familyId)
	defer dbResponse.Close()
	if err != nil {
		return nil, err
	}
	var memberList []*models.Member
	for i := 0; dbResponse.Next(); i++ {
		err = dbResponse.Scan(memberList[i])
		if err != nil {
			return nil, err
		}
	}
	return memberList, nil
}

func (fr *familyRepository) GetRootFolder(familyId int64) (*models.Folder, error) {
	queryResponse, err :=
		fr.db.Query(`select id,name,family_id from folders where family_id=$1 and name=$2`,
			familyId, "root")
	if err != nil {
		return nil, err
	}
	var folder models.Folder
	if queryResponse.Next() {
		err = queryResponse.Scan(&folder.Id, &folder.Name, &folder.FamilyId)
		if err != nil {
			return nil, err
		}
	}
	return &folder, err
}

func (fr *familyRepository) JoinFamily(familyId int64, user *models.User, password string) (bool, error) {
	ctx := context2.Background()
	txn, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer txn.Rollback()
	queryResponse, err :=
		txn.Query(`select families.sha_hash,folders.id from families 
                         inner join folders on folders.family_id=families.id 
                         where families.id=$1 and folders.name=$2`, familyId, "root")
	if err != nil {
		return false, err
	}
	var shaHash string
	var rootFolderId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&shaHash, &rootFolderId)
		if err != nil {
			return false, err
		}
	}
	err = queryResponse.Close()
	if err != nil {
		return false, err
	}
	if !matchPassword(password, shaHash) { // password check
		return false, errors.New("password is wrong")
	}
	queryResponse, err = txn.Query(`insert into folders(name,family_id) 
                                  values($1,$2) returning id`, user.Name+"'s Folder", familyId)
	if err != nil {
		return false, err
	}
	var userFolderId int64
	if queryResponse.Next() {
		err = queryResponse.Scan(&userFolderId)
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("cannot create user folder")
	}
	err = queryResponse.Close()
	if err != nil {
		return false, err
	}
	_, err =
		txn.Exec(`insert into folder_to_folder_map(parent_folder_id,folder_id) 
                         values($1,$2)`, rootFolderId, userFolderId)
	if err != nil {
		return false, err
	}
	_, err = txn.Exec(`insert into user_family_map(user_id,family_id) 
                                 values($1,$2)`, user.Id, familyId)
	if err != nil {
		return false, err
	}
	err = txn.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func matchPassword(password string, shaHash string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	if hash == shaHash {
		return true
	}
	return false
}
