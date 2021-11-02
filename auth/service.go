// JWT generate
package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct {
	
}

func ServiceBaru() *jwtService {
	return&jwtService{ }
}

var SECRET_KEY = []byte("Example_secret_key")

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_ID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, error := token.SignedString(SECRET_KEY)
	if error != nil {
		return signedToken, error
	}

	return signedToken, nil
}


// Vaidasi token
func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, error := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // HS256 adalah salah satu bentuk dari HMAC
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if error != nil {
		return token, error
	}

	return token, nil

}