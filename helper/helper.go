package helper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func APIResponse(message string, code int, status string, data interface{}) response {

	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonresponse := response{
		Meta: meta,
		Data: data,
	}
	fmt.Println(jsonresponse, "jsonresponse")
	return jsonresponse
}

func ErrorBadRequest(err error, c *gin.Context) {
	responseuser := APIResponse("error requet", http.StatusBadRequest, "error", nil)
	if err != nil {
		panic(err)

	}
	c.JSON(http.StatusBadRequest, responseuser)
}

func FormatValidationError(err error) []string {

	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
