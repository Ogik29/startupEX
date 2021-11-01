package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/startupex?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.RepositoryBaru(db)
	userService := user.ServiceBaru(userRepository)
	authService := auth.ServiceBaru()
	userHandler := handler.HandlerBaru(userService, authService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.Registeruser) //Register
	api.POST("/sessions", userHandler.Login) //Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailable) //Check email
	api.POST("/avatars", userHandler.UploadAvatar) //avatar

	router.Run()

}

