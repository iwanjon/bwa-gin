package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error create transaction", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	currentUser := c.MustGet("curresntuser").(user.User)
	input.User = currentUser

	trans, err := h.service.CreateTransaction(input)
	if err != nil {

		responseuser := helper.APIResponse("error create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	responsetrans := transaction.FormatTransaction(trans)
	responseuser := helper.APIResponse("success", http.StatusOK, "success", responsetrans)
	c.JSON(http.StatusOK, responseuser)

}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("curresntuser").(user.User)
	UserID := currentUser.ID
	transactions, err := h.service.GetTransactionByUserID(UserID)
	fmt.Println(transactions, "transs", UserID, currentUser)
	if err != nil {

		responseuser := helper.APIResponse("error get campaign transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	responseuser := helper.APIResponse("OK", http.StatusOK, "OK", transaction.FormatUserTransactions(transactions))
	// responseuser := helper.APIResponse("OK", http.StatusOK, "OK", transactions)
	c.JSON(http.StatusOK, responseuser)

}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error get campaign transaction", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	currentUser := c.MustGet("curresntuser").(user.User)
	input.User = currentUser

	transact, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {

		responseuser := helper.APIResponse("error get campaign transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}
	responseuser := helper.APIResponse("OK", http.StatusOK, "OK", transaction.FormatCampaignTransactions(transact))
	c.JSON(http.StatusOK, responseuser)

}

func (h *transactionHandler) GetNotif(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errormessage := gin.H{"errors": errors}
		responseuser := helper.APIResponse("error procedd transaction notification", http.StatusBadRequest, "error", errormessage)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	err = h.service.ProcessPayment(input)

	if err != nil {

		responseuser := helper.APIResponse("error procedd transaction notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseuser)
		return
	}

	responseuser := helper.APIResponse("OK", http.StatusOK, "OK", input)
	c.JSON(http.StatusOK, responseuser)

}
