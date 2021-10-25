package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	Findbyemail(email string) (User, error)
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


// login endpoint
func (r *repository) Findbyemail(email string) (User, error) {
	var user User
	error := r.db.Where("email = ?", email).Find(&user).Error
	if error != nil {
		return user, error
	}
	return user, nil
}