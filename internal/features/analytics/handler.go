package analytics

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

// GetDashboard godoc
// @Summary Dashboard financeiro
// @Description Retorna resumo financeiro e analytics do mês
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param month query int false "Mês (1-12)"
// @Param year query int false "Ano (ex: 2026)"
// @Success 200 {object} DashboardDTO
// @Router /analytics/dashboard [get]
func (h *Handler) GetDashboard(c *gin.Context) {

	now := time.Now()

	month := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	year := c.DefaultQuery("year", strconv.Itoa(int(now.Year())))

	userID := c.GetUint64("userID")

	dashboard, err := h.service.GetDashboard(userID, month, year)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
