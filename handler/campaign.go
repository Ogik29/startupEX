package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
// service yang menentukan repository mana yang di panggil
// repository: FindAll, FIndByUserID
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// get campaigns (list campaign endpoint)
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) // Fungsi dari "strconv.Atoi" untuk mengconvert keluaran menjadi int

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIresponse("Gagal mendapat campaign", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Sukses mendapat campaign", http.StatusOK, "Sukses", campaign.FormatCampaigns(campaigns))
    c.JSON(http.StatusOK, response)
}


// Campaign detail API (memunculkan hanya satu campaign yang dimiliki user)
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// api/vi/campaigns/1(id)
	// handler: mapping id yang di url ke struct input => service, call formatter
	// service: inputnya struct input => menagkap id di url, call repository
	// repository: get campaign by ID

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIresponse("Failed to get detail of campaign", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaign, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIresponse("Failed to get detail of campaign (2)", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	response := helper.APIresponse("Campaign detail", http.StatusOK, "Sukses", campaign)
	c.JSON(http.StatusOK, response)
}
