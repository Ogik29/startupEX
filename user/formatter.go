package user

type Userformatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageURL   string `json:"image_url"`
}

func Formatuser(user User, Token string) Userformatter {
	formatter := Userformatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      Token,
		ImageURL:   user.AvatarFileName,
	}

	return formatter
}