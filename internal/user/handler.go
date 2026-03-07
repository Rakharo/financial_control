package user

import (
	"net/http"
	"strconv"

	"financial_control/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login godoc
// @Summary Login de usuário
// @Description Autentica o usuário e retorna JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	user, err := h.service.Login(req.Login, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	response := LoginResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Login: user.Login,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers godoc
// @Summary Lista usuários
// @Description Retorna todos os usuários
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} UserResponse
// @Router /user [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []UserResponse

	for _, u := range users {
		response = append(response, UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Login: u.Login,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary Buscar usuário por ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.service.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Register godoc
// @Summary Criar usuário
// @Description Cria um novo usuário no sistema
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} map[string]string
// @Router /register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	err := h.service.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

// UpdateUser godoc
// @Summary Atualizar usuário
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body User true "User data"
// @Success 200 {object} map[string]string
// @Router /user/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}
	err = h.service.UpdateUser(id, u.Name, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser godoc
// @Summary Deletar usuário
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Router /user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
