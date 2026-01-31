package handler

import (
	"backer/helper"
	"backer/transaction"
	"backer/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		validationErrors := helper.FormatValidationError(err)
		errorMessages := gin.H{"errors": validationErrors}

		response := helper.APIResponse(helper.MsgInvalidTransactionInput, http.StatusBadRequest, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		if errors.Is(err, transaction.ErrCampaignNotFound) {
			response := helper.APIResponse(helper.MsgCampaignNotFound, http.StatusNotFound, "error", errorMessage)
			c.JSON(http.StatusNotFound, response)
			return
		}

		if errors.Is(err, transaction.ErrNotAuthorized) {
			response := helper.APIResponse(helper.MsgNotAuthorizedToViewTransactions, http.StatusForbidden, "error", errorMessage)
			c.JSON(http.StatusForbidden, response)
			return
		}

		response := helper.APIResponse(helper.MsgFailedToGetCampaignTransactions, http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(helper.MsgCampaignTransactionsRetrievedSuccess, http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse(helper.MsgFailedToGetUserTransactions, http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(helper.MsgUserTransactionsRetrievedSuccess, http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}
