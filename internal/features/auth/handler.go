package auth

import (
	"financial_control/internal/features/user"
	"financial_control/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *Service
	userService *user.Service
}

func NewHandler(authService *Service, userService *user.Service) *Handler {
	return &Handler{authService: authService, userService: userService}
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
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Informaçoes inválidas", err))
		return
	}

	var (
		u            *user.User
		accessToken  string
		refreshToken string
		err          error
	)

	if req.GoogleToken != "" {
		u, accessToken, refreshToken, err = h.authService.LoginWithGoogle(req.GoogleToken)
	} else {
		u, accessToken, refreshToken, err = h.authService.Login(req.Email, req.Password)
	}

	if err != nil {
		c.Error(err)
		return
	}

	providers, err := h.authService.GetUserProviders(u.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		User: user.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
		},
		Providers: providers,
	})
}

// LinkGoogle godoc
// @Summary Link google
// @Description Linka conta de usuário padrão ao google
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Google token"
// @Success 200 {object} map[string]string
// @Router /auth/link-google [post]
func (h *Handler) LinkGoogle(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req struct {
		GoogleToken string `json:"googleToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Informaçoes inválidas", err))
		return
	}

	err := h.authService.LinkGoogleAccount(userID, req.GoogleToken)
	if err != nil {
		c.Error(utils.NewBadRequest("Erro ao vincular conta. Email do Google diferente do usuário logado", err))
		return
	}

	c.JSON(200, gin.H{"message": "Conta Google vinculada"})
}

// RefreshToken godoc
// @Summary Refresh do token do usuario
// @Description atualiza o token do usuario com nova data de expieração
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Informações inválidas"})
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900, // 15 minutos
	})
}

// UpdateUserPassword godoc
// @Summary Alteraçao de senha
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordRequest true "Password data"
// @Success 200 {object} map[string]string
// @Router /auth/password [put]
func (h *Handler) UpdateUserPassword(c *gin.Context) {

	userID := c.GetUint64("userID")

	var req PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.authService.UpdateUserPassword(userID, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada"})
}

// Logout godoc
// @Summary logout de usuário
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid body"})
		return
	}

	userID := c.GetUint64("userID")
	err := h.authService.Logout(userID, req.RefreshToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Logout realizado"})
}
