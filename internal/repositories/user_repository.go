package repositories

import (
	"database/sql"
	"famdocs-zoo/models"
)

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
	GetUserFromId(int64) (*models.User, error)
	GetUserFamilies(*models.User) ([]*models.Family, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserFromId(userId int64) (*models.User, error) {
	user := new(models.User)
	dbResponse, err := ur.db.Query(`select id,name,phone,avatar from users where id=$1`, userId)
	defer dbResponse.Close()
	if err != nil {
		return nil, err
	}
	if dbResponse.Next() {
		err = dbResponse.Scan(&user.Id, &user.Name, &user.Phone, &user.Avatar)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (ur *userRepository) GetUserFamilies(user *models.User) ([]*models.Family, error) {
	dbResponse, err := ur.db.Query(
		`select f.id as id,f.name as name,f.owner as owner,
        f.sha_hash as sha_hash, folders.name as folder_name,
        folders.id as folder_id
        from user_family_map uf inner join families f
		on uf.family_id=f.id inner join folders on folders.family_id=f.id 
        where uf.user_id=$1 and folders.name='root'`, user.Id)
	defer dbResponse.Close()
	if err != nil {
		return nil, err
	}
	var familyList []*models.Family
	for i := 0; dbResponse.Next(); i++ {
		family := new(models.Family)
		err = dbResponse.Scan(&family.Id, &family.Name,
			&family.Owner, &family.ShaHash, &family.Root.Name, &family.Root.Id)
		family.Root.FamilyId = family.Id
		family.Root.ParentFolderId = -1
		if err != nil {
			return nil, err
		}
		familyList = append(familyList, family)
	}
	return familyList, nil
}
