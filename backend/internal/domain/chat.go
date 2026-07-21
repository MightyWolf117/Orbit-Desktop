package domain

import "context"

// ChatMessage representa un mensaje individual en la conversación
type ChatMessage struct {
	Role    string `json:"role"`    // "user" o "model"
	Content string `json:"content"` // Contenido del mensaje
}

// ChatRequest representa la solicitud entrante al endpoint de chat
type ChatRequest struct {
	Messages      []ChatMessage `json:"messages"`
	PersonalityID *int          `json:"personality_id,omitempty"`
	GenerateTitle bool          `json:"generate_title,omitempty"`
	ChatCode      *int64        `json:"chat_code,omitempty"` // Código único del chat desde el frontend
	Model         string        `json:"model,omitempty"`
	Temperature   *float32      `json:"temperature,omitempty"`
}

// ChatResponse representa la respuesta que devolvemos al frontend
type ChatResponse struct {
	Response  string     `json:"response"`
	Title     string     `json:"title,omitempty"` // Título generado para el historial, si GenerateTitle era true
	Historial *Historial `json:"historial,omitempty"` // Objeto historial guardado en DB
	Error     string     `json:"error,omitempty"`
}

// AIService define los métodos de negocio para interactuar con la IA
type AIService interface {
	GenerateResponse(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}
