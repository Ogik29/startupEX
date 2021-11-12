package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	Findbyemail(email string) (User, error)
	FindbyID(ID int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func RepositoryBaru(db *gorm.DB) *repository {
	return &repository{db}
}


// Register endpoint
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}


// login endpoint & check email
func (r *repository) Findbyemail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}


// Avatar endpoint
func (r *repository) FindbyID(ID int) (User, error) {
	var user User
	err := r.db.Where("id = ?", ID).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}