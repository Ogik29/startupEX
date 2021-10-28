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
	PasswordHash, error := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if error != nil {
		return user, error
	}
	user.PasswordHash = string(PasswordHash)
	user.Role = "user"

	newuser, error := s.repository.Save(user)
	if error != nil {
		return newuser, error
	}

	return newuser, nil

}


// login endpoint
func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, error := s.repository.Findbyemail(email)
	if error != nil {
		return user, error
	}
	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}
    error = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if error != nil {
		return user, error
	}

	return user, nil

}


// Chech email endpoint
func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, error := s.repository.Findbyemail(email)
	if error != nil {
		return false, error
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

	user, error := s.repository.FindbyID(ID)
	if error != nil {
		return user, error
	}

	user.AvatarFileName = fileLocation

	updatedUser, error := s.repository.Update(user)
	if error != nil {
		return updatedUser, error
	}

	return updatedUser, nil
}