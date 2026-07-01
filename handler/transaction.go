package handler

import (
	"backer/helper"
	"backer/transaction"
	"backer/user"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

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

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		validationErrors := helper.FormatValidationError(err)
		errorMessages := gin.H{"errors": validationErrors}

		response := helper.APIResponse(helper.MsgInvalidTransactionInput, http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		if errors.Is(err, transaction.ErrCampaignNotFound) {
			response := helper.APIResponse(helper.MsgCampaignNotFound, http.StatusNotFound, "error", errorMessage)
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponse(helper.MsgFailedToCreateTransaction, http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(helper.MsgTransactionCreatedSuccessfully, http.StatusCreated, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusCreated, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	// Read the raw body first for logging, before binding it to the struct
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	log.Println("=== MIDTRANS NOTIFICATION RAW BODY ===")
	log.Println(string(bodyBytes))
	log.Println("=== CONTENT-TYPE ===", c.GetHeader("Content-Type"))

	// Restore the body so it can be read again by ShouldBindJSON
	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("=== BIND ERROR ===", err.Error())
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Printf("=== PARSED INPUT === %+v\n", input)

	err = h.service.ProcessPayment(input)
	if err != nil {
		log.Println("=== PROCESS PAYMENT ERROR ===", err.Error())

		if errors.Is(err, transaction.ErrInvalidSignature) {
			response := helper.APIResponse("Invalid signature", http.StatusUnauthorized, "error", nil)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		if errors.Is(err, transaction.ErrTransactionNotFound) {
			response := helper.APIResponse("Transaction not found", http.StatusNotFound, "error", nil)
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
