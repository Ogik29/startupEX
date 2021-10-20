package user

import "golang.org/x/crypto/bcrypt"

type Service interface {
	Registeruser(input Register) (User, error)
}

type service struct {
	repository Repository
}

func ServiceBaru(repository Repository) *service {
	return &service{repository}
}

func (s *service) Registeruser(input Register) (User, error) {
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
