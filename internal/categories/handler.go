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

// GetAll godoc
// @Summary Listar categorias
// @Description Retorna todas as categorias disponíveis para o usuário (categorias próprias + categorias padrão do sistema)
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Category
// @Failure 500 {object} map[string]string
// @Router /category [get]
func (h *Handler) GetAll(c *gin.Context) {

	userID := c.GetUint64("user_id")

	categories, err := h.service.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetByID godoc
// @Summary Buscar categoria
// @Description Retorna uma categoria pelo ID
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/{id} [get]
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

// Create godoc
// @Summary Criar categoria
// @Description Cria uma nova categoria personalizada para o usuário
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body CategoryRequest true "Dados da categoria"
// @Success 201 {object} Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category [post]
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

// Update godoc
// @Summary Atualizar categoria
// @Description Atualiza uma categoria criada pelo usuário
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Param id path int true "ID da categoria"
// @Param category body CategoryRequest true "Dados atualizados da categoria"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /category/{id} [put]
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

// Delete godoc
// @Summary Deletar categoria
// @Description Remove uma categoria criada pelo usuário
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "ID da categoria"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /category/{id} [delete]
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
