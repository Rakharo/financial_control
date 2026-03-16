package transaction

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetTransactions godoc
// @Summary Lista transações
// @Description Retorna todas as transações
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Param page query int false "1"
// @Param limit query int false "10"
// @Param month query int false "Mês (1-12)"
// @Param year query int false "Ano (ex: 2026)"
// @Success 200 {array} TransactionResponse
// @Router /transaction [get]
func (h *Handler) GetTransactions(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	now := time.Now()

	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(int(now.Year()))))

	userID := c.GetUint64("userID")
	transactions, total, err := h.service.GetAllTransactions(userID, page, limit, month, year)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"page":         page,
		"limit":        limit,
		"total":        total,
		"month":        month,
		"year":         year,
	})
}

// GetTransactionByID godoc
// @Summary Buscar transação por ID
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} TransactionResponse
// @Failure 404 {object} map[string]string
// @Router /transaction/{id} [get]
func (h *Handler) GetTransactionByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIDInterface, _ := c.Get("userID")
	userID := uint64(userIDInterface.(int64))

	transaction, err := h.service.GetByID(id, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// CreateTransaction godoc
// @Summary Criar transação
// @Description Cria uma nova transação no sistema
// @Tags Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body TransactionRequest true "Transaction data"
// @Success 201 {object} map[string]string
// @Router /transaction [post]
func (h *Handler) CreateTransaction(c *gin.Context) {

	var dto TransactionRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint64("userID")

	transaction, err := h.service.Create(userID, dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// UpdateTransaction godoc
// @Summary Atualizar transação
// @Tags Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body TransactionRequest true "Transaction data"
// @Success 200 {object} map[string]string
// @Router /transaction/{id} [put]
func (h *Handler) UpdateTransaction(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto TransactionRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	userID := c.GetUint64("userID")

	transaction, err := h.service.Update(id, userID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary Deletar transação
// @Tags Transactions
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Router /transaction/{id} [delete]
func (h *Handler) DeleteTransaction(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.Delete(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
