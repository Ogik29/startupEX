package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Registeruser(input RegisterInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func ServiceBaru(repository Repository) *service {
	return &service{repository}
}


// Register endpoint
func (s *service) Registeruser(input RegisterInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(PasswordHash)
	user.Role = "user"

	newuser, err := s.repository.Save(user)
	if err != nil {
		return newuser, err
	}

	return newuser, nil

}


// login endpoint
func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.Findbyemail(email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil

}


// Chech email endpoint
func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.repository.Findbyemail(email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}
	return false, nil
}


// Avatar endpoint
func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// dapatkan user melalui ID
	// update atribut avatar file name
	// simpan perubahan avatar file name

	user, err := s.repository.FindbyID(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}


// Service untuk middleware avatar
func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindbyID(ID) 
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found on that ID")
	}

	return user, nil
	
} 