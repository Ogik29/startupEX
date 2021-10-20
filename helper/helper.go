package helper

//harus ditambahkan v10
import "github.com/go-playground/validator/v10" 

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIresponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	JSONresponse := Response{
		Meta: meta,
		Data: data,
	}

	return JSONresponse
}

func FormatValidationErrors(error error) []string {
	var errors []string

	for _, e := range error.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}