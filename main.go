package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.RepositoryBaru(db)
	userService := user.ServiceBaru(userRepository)
	userHandler := handler.Handlerbaru(userService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.Registeruser) //register
	api.POST("/sessions", userHandler.Login) //Login

	router.Run()

	
	
	// service
	// userinput := user.RegisterInput{}
	// userinput.Name = "Venti"
	// userinput.Occupation = "Archon"
	// userinput.Email = "bartobas@gmail.com"
	// userinput.Password = "sandi"

	// userService.Registeruser(userinput)
	

	// Repository
	// user := user.User {
	// 	ID: 3,
	// 	Name: "Ei",
	// 	Occupation: "Archon",
	// 	Email: "Ei@gmail.com",
	// 	PasswordHash: "test2",
	// 	AvatarFileName: "Shogun.jpg",
	// 	Role: "user",
	// }

	// userRepository.Save(user)

}

