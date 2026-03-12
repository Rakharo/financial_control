package auth

import (
	"financial_control/internal/features/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{userService: userService}
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
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	u, err := h.userService.Login(req.Login, req.Password)

	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, _ := GenerateAccessToken(uint64(u.ID))
	refreshToken, _ := GenerateRefreshToken(uint64(u.ID))

	c.JSON(200, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		User: user.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Login: u.Login,
		},
	})
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
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	claims, err := ValidateToken(req.RefreshToken)

	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	newToken, _ := GenerateAccessToken(claims.UserID)

	c.JSON(200, gin.H{
		"access_token": newToken,
		"expires_in":   900,
	})
}
