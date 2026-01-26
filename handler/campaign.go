package handler

import (
	"backer/campaign"
	"backer/helper"
	"backer/user"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns retrieved successfully", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Invalid campaign ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Campaign not found", http.StatusNotFound, "error", errorMessage)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Campaign detail retrieved successfully", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Invalid input", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Failed to create campaign", http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Campaign created successfully", http.StatusCreated, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusCreated, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Invalid campaign ID", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Invalid input", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		if strings.Contains(err.Error(), "not found") {
			response := helper.APIResponse("Campaign not found", http.StatusNotFound, "error", errorMessage)
			c.JSON(http.StatusNotFound, response)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "not authorized") {
			response := helper.APIResponse("You are not authorized to update this campaign", http.StatusForbidden, "error", errorMessage)
			c.JSON(http.StatusForbidden, response)
			return
		}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign updated successfully", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Invalid input", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("No file uploaded", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Invalid file type. Only JPG, JPEG, and PNG are allowed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	maxSize := int64(2 << 20)
	if file.Size > maxSize {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("File size too large. Maximum 2MB", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d-%d-%s", userID, timestamp, file.Filename)
	path := fmt.Sprintf("images/%s", fileName)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to save file to server", http.StatusInternalServerError, "error", data)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error(), "is_uploaded": false}

		if strings.Contains(err.Error(), "not found") {
			response := helper.APIResponse("Campaign not found", http.StatusNotFound, "error", errorMessage)
			c.JSON(http.StatusNotFound, response)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") {
			response := helper.APIResponse("You are not authorized to upload image for this campaign", http.StatusForbidden, "error", errorMessage)
			c.JSON(http.StatusForbidden, response)
			return
		}

		response := helper.APIResponse("Failed to save image to database", http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image uploaded successfully", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
