package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func Handlerbaru(userService user.Service) *userHandler{
	return &userHandler{userService}
}

func (h *userHandler) Registeruser(c *gin.Context) {
	// tangkap input dari user
	// map inputan dari user ke struct Register
	// struct di user di passing sebagai parameter service

	var input user.RegisterInput

	error := c.ShouldBindJSON(&input)
	if error != nil{
		errors := helper.FormatValidationErrors(error)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Akun gagal terbuat", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newuser, error := h.userService.Registeruser(input)

	if error != nil{
		response := helper.APIresponse("Akun gagal terbuat", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.Formatuser(newuser, "token(tumbal :v)")

	response := helper.APIresponse("Akun sudah terbuat", http.StatusOK, "Sukses", formatter)
	
	c.JSON(http.StatusOK, response)

}


// login endpoint
func (h *userHandler) Login(c *gin.Context) {
    // user memasukkan input (email & password)
    // input ditangkap handler
    // mapping dari input user ke input struct
    // input struct passing ke service
    // di service mencari dengan bantuan repository user dengan email x
    // mencocokkan password 

	var input user.LoginInput

	error := c.ShouldBindJSON(&input)
	if error != nil {
		errors := helper.FormatValidationErrors(error)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Login gagal", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinuser, error := h.userService.Login(input)

	if error != nil {
		errormessage := gin.H{"errors": error.Error()}

		response := helper.APIresponse("Login gagal", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.Formatuser(loggedinuser, "token(tumbal :v)")

	response := helper.APIresponse("Login sukses", http.StatusOK, "Sukses", formatter)
	
	c.JSON(http.StatusOK, response)
	
}