package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func RepositoryBaru(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	error := r.db.Create(&user).Error
	if error != nil {
		return user, error
	}
	return user, nil
}