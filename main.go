package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/transaction"
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
    
	// repository
	userRepository := user.RepositoryBaru(db)
	campaignRepository := campaign.RepositoryBaru(db)
	transactionRepository := transaction.RepositoryBaru(db)
    
	// service
	userService := user.ServiceBaru(userRepository)
	campaignService := campaign.ServiceBaru(campaignRepository)
	authService := auth.ServiceBaru()
	transactionService := transaction.ServiceBaru(transactionRepository, campaignRepository)

    // handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Static("/images", "./images") // berfungsi untuk menampilkan gambar di postman

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.Registeruser)                                             // Register
	api.POST("/sessions", userHandler.Login)                                                 // Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailable)                             // Check email
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) // avatar

	api.GET("/campaigns", campaignHandler.GetCampaigns)                                                 // Get Campaigns
	api.GET("/campaign/:id", campaignHandler.GetCampaign)                                               // Campaign detail
	api.POST("/campaign", authMiddleware(authService, userService), campaignHandler.CreateCampaign)     // Create campaign
	api.PUT("/campaign/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)  // Update campaign
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage) // Upload campaign image

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions) // Get campaign transaction
    api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions) // Get user transaction

	// Fungsi dari "authMiddleware(authService, userService)" untuk mengetahui siapa user
	// yang melakukan request seperti mengupload avatar/membuat campaign, dll.

	router.Run()

}

// Middleware
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_ID"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIresponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
