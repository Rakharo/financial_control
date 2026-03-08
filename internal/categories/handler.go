package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {

	userID := c.GetUint64("user_id")

	categories, err := h.service.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) GetByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	category, err := h.service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "categoria não encontrada"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) Create(c *gin.Context) {

	userID := c.GetUint64("user_id")

	var dto CategoryRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	category, err := h.service.Create(userID, dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handler) Update(c *gin.Context) {

	userID := c.GetUint64("user_id")

	idParam := c.Param("id")

	id, _ := strconv.ParseUint(idParam, 10, 64)

	var dto CategoryRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.service.Update(id, userID, dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {

	userID := c.GetUint64("user_id")

	idParam := c.Param("id")

	id, _ := strconv.ParseUint(idParam, 10, 64)

	err := h.service.Delete(id, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}
