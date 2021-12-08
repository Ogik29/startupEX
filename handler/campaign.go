package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
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

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIresponse("Failed to get detail of campaign (2)", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	response := helper.APIresponse("Campaign detail", http.StatusOK, "Sukses", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}


// Create campaign endpoint
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// tangkap parameter dari user ke input struct
    // ambil current user dari jwt/handler
    // panggil service, parameter input struct (dan juga buat slug)
    // panggil repository untuk simpan data campaign baru

	var input campaign.CreateCampaignInput 
	
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Gagal membuat campaign", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
    
	// berfungsi agar tau siapa user yang sedang membuat campaign
	currentUser := c.MustGet("currentUser").(user.User)
    input.User = currentUser
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIresponse("Gagal membuat campaign", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Berhasil membuat campaign", http.StatusOK, "Sukses", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
	
}


// Update campaign endpoint
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	// user memasukkan input
	// handler
	// mapping dari input ke input struct (ada 2 dari user dan dari uri)
	// input dari user, dan juga input dari uri (passing ke service)
	// service (find campaign by ID, tangkap parameter)
	// repository update data campaign

	var inputID campaign.GetCampaignDetailInput // untuk mendapatkan id campaign 

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIresponse("Failed to update campaign", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput // untuk mendapatkan data campaign yang akan di update
	
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errormessage := gin.H{"errors": errors}

		response := helper.APIresponse("Failed to update campaign", http.StatusUnprocessableEntity, "Eror", errormessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// berfungsi agar tau siapa user yang sedang mengupdate campaign
	currentUser := c.MustGet("currentUser").(user.User)
    inputData.User = currentUser
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIresponse("Failed to update campaign", http.StatusBadRequest, "Eror", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Success to update campaign", http.StatusOK, "Sukses", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}
