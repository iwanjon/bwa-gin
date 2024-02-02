package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authservice auth.Service
}

func NewHandler(userService user.Service, authservice auth.Service) *userHandler {
	return &userHandler{userService, authservice}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUser

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error requet", http.StatusUnprocessableEntity, "error", errormessage)
		c.JSON(http.StatusUnprocessableEntity, responseuser)
		return
	}
	// helper.ErrorBadRequest(err, c)
	newuser, err := h.userService.RegisterUser(input)
	if err != nil {
		responseuser := helper.APIResponse("error requet", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	token, err := h.authservice.GenerateToken(newuser.ID)
	if err != nil {
		responseuser := helper.APIResponse("error token", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	fmt.Println(newuser, err)
	// helper.ErrorBadRequest(err, c)
	newformat := user.Formatuser(newuser, token)

	responseuser := helper.APIResponse("success cretaed user", http.StatusOK, "ok", newformat)
	fmt.Println(newformat, responseuser, "responseuser")
	c.JSON(http.StatusOK, responseuser)
}

func (h *userHandler) LoginUser(c *gin.Context) {
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error login", http.StatusUnprocessableEntity, "error", errormessage)
		c.JSON(http.StatusUnprocessableEntity, responseuser)
		return
	}
	fmt.Println(input, input, "input")
	loginuser, err := h.userService.LoginUser(input)
	fmt.Println(input, input, "input")
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error login", http.StatusUnprocessableEntity, "error", errormessage)
		c.JSON(http.StatusUnprocessableEntity, responseuser)
		return
	}

	token, err := h.authservice.GenerateToken(loginuser.ID)
	if err != nil {
		responseuser := helper.APIResponse("login failed ", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	formater := user.Formatuser(loginuser, token)
	respon := helper.APIResponse("sukses login", http.StatusOK, "sukses", formater)
	c.JSON(http.StatusOK, respon)
}

func (h *userHandler) HandlerEmailAvailability(c *gin.Context) {

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error login", http.StatusUnprocessableEntity, "error", errormessage)
		c.JSON(http.StatusUnprocessableEntity, responseuser)
		return
	}
	isemailavailable, err := h.userService.CheckEmailAvailable(input)

	if err != nil {

		errormessage := gin.H{"errors": "server errror"}
		responseuser := helper.APIResponse("failed to check email", http.StatusBadGateway, "error", errormessage)
		c.JSON(http.StatusBadGateway, responseuser)
		return
	}

	data := gin.H{
		"is_available": isemailavailable,
	}
	metaMessage := "email has been registered"

	if isemailavailable {
		metaMessage = "email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("uploaded failed", http.StatusBadRequest, "failed", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userid := c.MustGet("curresntuser").(user.User)
	userID := userid.ID
	// userID := 1
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	// path := "images/" + file.Filename
	// fmt.Println(path, file, "path and file")
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("uploaded failed 1", http.StatusBadRequest, "failed", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("uploaded 2 failed", http.StatusBadRequest, "failed", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("uploaded 3 success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) FetchUser(c *gin.Context) {

	currentUser := c.MustGet("curresntuser").(user.User)

	formatter := user.Formatuser(currentUser, "")

	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}
