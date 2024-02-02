package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService: campaignService, userService: userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campains, err := h.campaignService.FindCampaigns(0)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campains})
}

func (h *campaignHandler) CreateView(c *gin.Context) {
	userList, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	input := campaign.FormCreateCampaignInput{}
	input.Users = userList

	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.FormCreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		users, er := h.userService.GetAllUsers()
		if er != nil {
			c.HTML(http.StatusBadRequest, "error.html", nil)
			return
		}
		log.Println(users, err, "this", er)
		input := campaign.FormCreateCampaignInput{}
		input.Users = users
		input.Error = err
		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	userCurrent, err := h.userService.GetUserById(input.UserID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	inputCampaign := campaign.CreateCampaignInput{}
	inputCampaign.Name = input.Name
	inputCampaign.Description = input.Description
	inputCampaign.GoalAmount = input.GoalAmount
	inputCampaign.ShortDescription = input.ShortDescription
	inputCampaign.Perks = input.Perks
	inputCampaign.User = userCurrent

	_, err = h.campaignService.CreateCampaign(inputCampaign)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	// campains, err := h.campaignService.FindCampaigns(0)
	c.Redirect(http.StatusOK, "campaign_index.html")
	// return
}

func (h *campaignHandler) UploadImageView(c *gin.Context) {
	strId := c.Param("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": id})
}

func (h *campaignHandler) UploadImage(c *gin.Context) {

	strId := c.Param("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	camID := id
	// userID := 1
	path := fmt.Sprintf("images/campaign-%d-%s", camID, file.Filename)
	// path := "images/" + file.Filename
	// fmt.Println(path, file, "path and file")
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	// h.campaignService.
	camIdStruct := campaign.GetCampaignDetailInput{ID: camID}

	currentCampaign, err := h.campaignService.GetCampaignById(camIdStruct)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	campaignImageInput := campaign.CreateCampaignImageInput{
		CampaignID: camID,
		IsPrimary:  true,
		User:       currentCampaign.User,
	}
	_, err = h.campaignService.SaveCampaignImage(campaignImageInput, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) EditCampaign(c *gin.Context) {
	strID := c.Param("id")

	intID, err := strconv.Atoi(strID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	detailInput := campaign.GetCampaignDetailInput{ID: intID}

	curentCampaign, err := h.campaignService.GetCampaignById(detailInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		// fmt.Println("sini", curentCampaign, detailInput)
		// c.HTML(http.StatusOK, "campaign_edit.html", gin.H{"Campaign": curentCampaign, "Error": err})
		return
	}
	// var ac = make(map[string]interface{})

	// as := gin.H(ac)
	input := campaign.FormUpdateCampaignInput{}
	input.Description = curentCampaign.Description
	input.ID = intID
	input.GoalAmount = curentCampaign.GoalAmount
	input.Name = curentCampaign.Name
	input.Perks = curentCampaign.Perks
	input.ShortDescription = curentCampaign.ShortDescription
	input.User = curentCampaign.User

	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) PostUpdateCampaign(c *gin.Context) {
	strID := c.Param("id")

	intID, err := strconv.Atoi(strID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	dataToInout := campaign.FormUpdateCampaignInput{}
	// data2 := campaign.FormUPdateCampaignInput{}
	err = c.ShouldBind(&dataToInout)
	fmt.Println(dataToInout, "lolo", err, "koko")
	// err = c.ShouldBind(&data2)
	// fmt.Println(data2, "koa", err, "koko2")

	if err != nil {
		dataToInout.Error = err
		dataToInout.ID = intID
		c.HTML(http.StatusOK, "campaign_edit.html", dataToInout)
		return
	}
	detailInput := campaign.GetCampaignDetailInput{ID: intID}

	curentCampaign, err := h.campaignService.GetCampaignById(detailInput)
	// input := campaign.FormUpdateCampaignInput{}
	// input.Description = curentCampaign.Description
	// input.ID = intID
	// input.GoalAmount = curentCampaign.GoalAmount
	// input.Name = curentCampaign.Name
	// input.Perks = curentCampaign.Perks
	// input.ShortDescription = curentCampaign.ShortDescription
	// input.User = curentCampaign.User
	if err != nil {
		// c.HTML(http.StatusInternalServerError, "error.html", nil)
		fmt.Println("sini", curentCampaign, detailInput)
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		// c.HTML(http.StatusOK, "campaign_edit.html", input)
		return
	}
	dataToSAve := campaign.CreateCampaignInput{
		Name:             dataToInout.Name,
		ShortDescription: dataToInout.ShortDescription,
		Description:      dataToInout.Description,
		GoalAmount:       dataToInout.GoalAmount,
		Perks:            dataToInout.Perks,
		User:             curentCampaign.User,
	}
	fmt.Println(dataToSAve, "Dfffd")
	_, err = h.campaignService.UpdateCampaign(detailInput, dataToSAve)
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)

		return
	}
	fmt.Println(intID)
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) ShowDetailCampaign(c *gin.Context) {

	str_id := c.Param("id")
	int_id, err := strconv.Atoi(str_id)
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)

		return
	}
	current_campaign := campaign.GetCampaignDetailInput{ID: int_id}
	campain, err := h.campaignService.GetCampaignById(current_campaign)
	fmt.Println(campain.CurrentAmountFormatIDR())
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)

		return
	}
	jeson, err := json.Marshal(campain)
	if err != nil {
		log.Println("error json MArshal")
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	var mapCampaign = make(map[string]interface{})
	err = json.Unmarshal(jeson, &mapCampaign)
	if err != nil {
		log.Println("error json UNMArshal")
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	mapCampaign["CurrentAmountFormatIDR"] = campain.CurrentAmountFormatIDR()
	mapCampaign["GoalAmountFormatIDR"] = campain.GoalAmountFormatIDR()
	ginH := gin.H(mapCampaign)
	c.HTML(http.StatusOK, "campaign_show.html", ginH)
}
