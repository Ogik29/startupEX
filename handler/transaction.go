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

// campaign transactions endpoint
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
