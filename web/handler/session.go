package handler

import (
	"bwastartup/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionHanlder struct {
	userService user.Service
}

func NewSessionHandler(userService user.Service) *sessionHanlder {
	return &sessionHanlder{userService}
}

func (h *sessionHanlder) NewSession(c *gin.Context) {

	c.HTML(http.StatusOK, "session_new.html", nil)
}

func (h *sessionHanlder) PostSession(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	uu, err := h.userService.LoginUser(input)
	if err != nil || uu.Role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session := sessions.Default(c)
	session.Set("userId", uu.ID)
	session.Save()
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *sessionHanlder) Destroy(c *gin.Context) {

	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}
