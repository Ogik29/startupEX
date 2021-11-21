package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	campaignRepository := campaign.RepositoryBaru(db)

	userService := user.ServiceBaru(userRepository)
	campaignService := campaign.ServiceBaru(campaignRepository)
	authService := auth.ServiceBaru()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images") // berfungsi untuk menampilkan gambar di postman

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.Registeruser) // Register
	api.POST("/sessions", userHandler.Login) // Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailable) // Check email
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) // avatar

	api.GET("/campaigns", campaignHandler.GetCampaigns) // Get Campaigns

	router.Run()

}


// Middleware
func authMiddleware (authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization") 
	
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
		
		// Bearer (spasi) token = yang diambil hanya tokennya saja
		tokenString := ""
		tokenArray := strings.Split(authHeader, " ")
		if len(tokenArray) == 2 {
			tokenString = tokenArray[1]
		}
        
		// Validasi token
	    token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
        
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		userID := int(claim["user_ID"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		c.Set("currentUser", user)
    }
	
}
   

