package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"orbit-backend/internal/domain"
)

type HistorialHandler struct {
	repo domain.HistorialRepository
}

func NewHistorialHandler(repo domain.HistorialRepository) *HistorialHandler {
	return &HistorialHandler{
		repo: repo,
	}
}

func (h *HistorialHandler) GetAll(c *gin.Context) {
	historiales, err := h.repo.GetHistorials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo historiales: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, historiales)
}

func (h *HistorialHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.repo.DeleteHistorial(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando historial: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
