package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}


// get campaign transactions endpoint
func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	// parameter di uri
    // tangkap parameter mapping ke input struct
    // panggil service, input struct sebagai parameterya
    // service, berbekal campaign id bisa panggil repository
    // repository mencari data transaction suatu campaign

	var input transaction.GetCampaignTransactionsInput
	
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIresponse("Failed to get campaign's transaction", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// berfungsi agar tau siapa user yang sedang melihat campaign's transaction (Authorization)
	currentUser := c.MustGet("currentUser").(user.User)
    input.User = currentUser
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIresponse("Failed to get campaign's transaction", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Campaign's transaction", http.StatusOK, "Sukses", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}


// get user transactions endpoint
func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	// handler
	// ambil nilai user dari jwt/middleware
	// service
	// repo => ambil data transactions (preload campaign)
    
	// Authorization
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIresponse("Failed to get user's transaction", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("User's transaction", http.StatusOK, "Sukses", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}


// user create transaction endpoint 
func (h *transactionHandler) CreateTransaction(c *gin.Context) {
    // ada input dari user
    // handler tangkap input terus di mapping ke input struct
    // panggil service buat transaksi, manggil sistem midtrans
    // panggil repository create new transaction data  

	var input transaction.CreateTransactionInput 

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Gagal membuat transaction", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

		// Authorization
		currentUser := c.MustGet("currentUser").(user.User)
		input.User = currentUser
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	
		newTransaction, err := h.service.CreateTransaction(input)
		if err != nil {
			response := helper.APIresponse("Gagal membuat transaction", http.StatusBadRequest, "Eror", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	
		response := helper.APIresponse("Berhasil membuat transaction", http.StatusOK, "Sukses", transaction.FormatTransaction(newTransaction))
		c.JSON(http.StatusOK, response)
}


// notification payment midtrans
func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIresponse("Failed to procces nitification", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
    
	err = h.service.ProccesPayment(input)
	if err != nil {
		response := helper.APIresponse("Failed to procces nitification", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input) 
	// dibuat sederhana dari yang lain karena yg mengakses endpoint ini bukanlah client atau semacamnya melainkan sistem midtrans 
}
