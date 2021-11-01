// JWT generate
package auth

import "github.com/dgrijalva/jwt-go"

type Service interface {
	GenerateToken(userID int, userName string, userEmail string) (string, error)
}

type jwtService struct {
	
}

func ServiceBaru() *jwtService {
	return&jwtService{ }
}

var SECRET_KEY = []byte("Example_secret_key")

func (s *jwtService) GenerateToken(userID int, userName string, userEmail string) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_ID"] = userID
	claim["user_name"] = userName
	claim["email"] = userEmail

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, error := token.SignedString(SECRET_KEY)
	if error != nil {
		return signedToken, error
	}

	return signedToken, nil
}