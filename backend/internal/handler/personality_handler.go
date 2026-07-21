package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"orbit-backend/internal/domain"
)

type PersonalityHandler struct {
	repo domain.PersonalityRepository
}

func NewPersonalityHandler(repo domain.PersonalityRepository) *PersonalityHandler {
	return &PersonalityHandler{
		repo: repo,
	}
}

func (h *PersonalityHandler) GetAll(c *gin.Context) {
	personalities, err := h.repo.GetAllPersonalities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener personalidades: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, personalities)
}

func (h *PersonalityHandler) Create(c *gin.Context) {
	var p domain.Personality
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := h.repo.CreatePersonality(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear personalidad: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Personalidad creada con éxito"})
}

func (h *PersonalityHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID es requerido"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero"})
		return
	}

	var p domain.Personality
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := h.repo.UpdatePersonality(id, &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar personalidad: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personalidad actualizada con éxito"})
}
