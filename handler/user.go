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


// Chech email endpoint : untuk mengetahui apakah email sudah terdaftar / belum
func (h *userHandler) CheckEmailAvailable(c *gin.Context) {
	// ada inputan email dari user
	// inputan email di mapping ke struct input
	// struct input dipassing ke service
	// service akan memanggil repository email sudah terdaftar atau belum
	// repository mengecek ke db 

	var input user.CheckEmailInput

	error := c.ShouldBindJSON(&input)
	if error != nil {
		errors := helper.FormatValidationErrors(error)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Checking email failed", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, error := h.userService.IsEmailAvailable(input)
	if error != nil {
		errormessage := gin.H{"errors": "Terjadi kesalahan"}
        response := helper.APIresponse("Checking email failed", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metamessage := "Email has been registered"
	if isEmailAvailable {
		metamessage = "Email is available"
	}

	response := helper.APIresponse(metamessage, http.StatusOK, "Sukses", data)
	c.JSON(http.StatusOK, response)

}