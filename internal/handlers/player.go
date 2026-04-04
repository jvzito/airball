package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jvzito/airball/internal/service"
)

type PlayerHandler struct{ svc *service.PlayerService }

func NewPlayerHandler(svc *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

func (h *PlayerHandler) GetLeaders(c *gin.Context) {
	category := strings.ToUpper(c.Param("category"))
	season := c.DefaultQuery("season", "2025-26")

	valid := map[string]bool{"PTS": true, "AST": true, "REB": true, "STL": true, "BLK": true, "FG3M": true}
	if !valid[category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "use: PTS, AST, REB, STL, BLK ou FG3M"})
		return
	}

	leaders, err := h.svc.GetLeaders(c.Request.Context(), season, category)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": leaders})
}

func (h *PlayerHandler) GetShotChart(c *gin.Context) {
	playerID := c.Param("id")
	season := c.DefaultQuery("season", "2025-26")

	shots, err := h.svc.GetShotChart(c.Request.Context(), playerID, season)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shots})
}
