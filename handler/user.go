package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func HandlerBaru(userService user.Service, authService auth.Service) *userHandler{
	return &userHandler{userService, authService}
}


// Register endpoint
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

	token, error := h.authService.GenerateToken(newuser.ID, newuser.Name, newuser.Email)
	if error != nil{
		response := helper.APIresponse("Akun gagal terbuat", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.Formatuser(newuser, token)

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

	token, error := h.authService.GenerateToken(loggedinuser.ID, loggedinuser.Name, loggedinuser.Email)
	if error != nil {
		errormessage := gin.H{"errors": error.Error()}

		response := helper.APIresponse("Login gagal", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.Formatuser(loggedinuser, token)

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


// Avatar (foto profil) endpoint
func (h *userHandler) UploadAvatar(c *gin.Context) {
	// input dari use 
	// simpan gambar di folder "images/"
	// di service memanggil repo
	// JWT (sementara hardcore, seakan-akan user yg login ID = 2)
	// repo ambill data user ID = 2
	// repo update data user, simpan lokasi file

	file, error := c.FormFile("avatar")
	if error != nil {
		data := gin.H{"Is_uploaded": false}
		response := helper.APIresponse("Failed to upload avatar image", http.StatusBadRequest, "Eror", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// sementara
	userID := 2
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	error = c.SaveUploadedFile(file, path)
	if error != nil {
		data := gin.H{"Is_uploaded": false}
		response := helper.APIresponse("Failed to upload avatar image", http.StatusBadRequest, "Eror", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
    
	// sementara
	_, error = h.userService.SaveAvatar(userID, path)
	if error != nil {
		data := gin.H{"Is_uploaded": false}
		response := helper.APIresponse("Failed to upload avatar image", http.StatusBadRequest, "Eror", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"Is_uploaded": true}
	response := helper.APIresponse("Avatar succesfuly uploaded", http.StatusOK, "Sukses", data)

	c.JSON(http.StatusOK, response)

}