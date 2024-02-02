package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {

		responseuser := helper.APIResponse("error updated campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	var inputdata campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputdata)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error update campaign", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	userid := c.MustGet("curresntuser").(user.User)
	inputdata.User = userid

	updated, err := h.service.UpdateCampaign(input, inputdata)
	if err != nil {
		responseuser := helper.APIResponse("error update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	responseuser := helper.APIResponse("success Update campaign", http.StatusOK, "success", campaign.FormatCampaign(updated))
	c.JSON(http.StatusOK, responseuser)

}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error create campaign", http.StatusBadRequest, "error create campaign", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	userid := c.MustGet("curresntuser").(user.User)
	input.User = userid

	newcampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		responseuser := helper.APIResponse("error create campaign", http.StatusBadRequest, "error create campaign", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	responseuser := helper.APIResponse("success create campaign", http.StatusOK, "success create campaign", campaign.FormatCampaign(newcampaign))
	c.JSON(http.StatusOK, responseuser)

}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigs, err := h.service.FindCampaigns(userId)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error requet", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	responseuser := helper.APIResponse("success", http.StatusOK, "success", campaign.FormatCampaigns(campaigs))
	c.JSON(http.StatusOK, responseuser)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {

		responseuser := helper.APIResponse("error INPUT", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	campain, err := h.service.GetCampaignById(input)
	if err != nil {

		responseuser := helper.APIResponse("error fet detail", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	responseuser := helper.APIResponse("success", http.StatusOK, "success", campaign.FormatCampaignDetail(campain))
	c.JSON(http.StatusOK, responseuser)

}

func (h *campaignHandler) UploadCampaignImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error upload campaign images", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_upload": false}
		response := helper.APIResponse("failed to upload campaign images", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("curresntuser").(user.User)
	input.User = currentUser
	userId := currentUser.ID

	path := fmt.Sprintf("images/campaign-%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_upload": false}
		response := helper.APIResponse("failed to upload campaign images", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_upload": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_upload": true}
	response := helper.APIResponse("success to upload campaign images", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
