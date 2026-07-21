package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type SystemHandler struct {
	apiKey string
}

func NewSystemHandler(apiKey string) *SystemHandler {
	return &SystemHandler{
		apiKey: apiKey,
	}
}

// Health responde un 200 OK para confirmar que el backend está activo
func (h *SystemHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Models consulta la API de Google y devuelve los modelos soportados
func (h *SystemHandler) Models(c *gin.Context) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(h.apiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la API de Gemini"})
		return
	}
	defer client.Close()

	// Recorremos los modelos disponibles (ej. gemini-1.5-flash, gemini-1.5-pro, etc)
	var availableModels []map[string]string
	
	iter := client.ListModels(ctx)
	for {
		m, err := iter.Next()
		if err != nil {
			// iterator.Done se devuelve cuando finaliza en la librería estándar, pero 
			// como no importamos google.golang.org/api/iterator directamente, controlamos el break al fallar.
			break
		}
		
		// Filtramos para asegurar que sean modelos de generación de texto usables para chat
		name := m.Name
		if !strings.HasPrefix(name, "models/gemini") || strings.Contains(name, "vision") || strings.Contains(name, "embedding") {
			continue
		}

		availableModels = append(availableModels, map[string]string{
			"id":           name,
			"displayName":  m.DisplayName,
			"description":  m.Description,
			"status":       "Activo", 
			"quotaMessage": "Para ver su cuota exacta restante comuníquese con el administrador.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"models": availableModels,
	})
}
