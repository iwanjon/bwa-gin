package handler

import (
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (u *userHandler) Index(c *gin.Context) {
	users, err := u.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (s *userHandler) CreateView(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)

}

func (u *userHandler) CreateUser(c *gin.Context) {
	var dataInput user.FromRegisterInput
	err := c.ShouldBind(&dataInput)
	if err != nil {
		dataInput.Error = err
		c.HTML(http.StatusOK, "user_new.html", dataInput)
		return
	}

	registeredData := user.RegisterUser{Name: dataInput.Name, Email: dataInput.Email, Password: dataInput.Password, Occupation: dataInput.Occupation}

	userregister, err := u.userService.RegisterUser(registeredData)

	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", nil)
		return
	}

	log.Println(userregister)
	// c.HTML(http.StatusOK, "user_new.html", nil)
	c.Redirect(http.StatusFound, "/users")
}

func (u *userHandler) UpdateUser(c *gin.Context) {
	strId := c.Param("id")
	id, _ := strconv.Atoi(strId)

	userr, err := u.userService.GetUserById(id)

	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userRegistered := user.FormUpdateInput{}
	userRegistered.ID = id
	userRegistered.Email = userr.Email
	userRegistered.Name = userr.Name
	userRegistered.Occupation = userr.Occupation

	c.HTML(http.StatusOK, "user_edit.html", userRegistered)

}

func (u *userHandler) PostUpdateUser(c *gin.Context) {
	strId := c.Param("id")
	id, _ := strconv.Atoi(strId)

	// userr, err := u.userService.GetUserById(id)
	var dataInput user.FormUpdateInput
	err := c.ShouldBind(&dataInput)
	if err != nil {
		// dataInput.Error = err
		dataInput.Error = err
		c.HTML(http.StatusOK, "user_edit.html", dataInput)
		// c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	dataInput.ID = id
	log.Println(dataInput, "all input ")
	_, err = u.userService.UodateUser(dataInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}

func (u *userHandler) UploadAvatarView(c *gin.Context) {
	strId := c.Param("id")
	id, _ := strconv.Atoi(strId)
	c.HTML(http.StatusOK, "user_avatar.html", gin.H{"ID": id})
}

func (u *userHandler) UploadAvatarSubmit(c *gin.Context) {
	idstr := c.Param("id")
	id, _ := strconv.Atoi(idstr)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userID := id
	// userID := 1
	path := fmt.Sprintf("images/user-%d-%s", userID, file.Filename)
	// path := "images/" + file.Filename
	// fmt.Println(path, file, "path and file")
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	_, err = u.userService.SaveAvatar(userID, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}
