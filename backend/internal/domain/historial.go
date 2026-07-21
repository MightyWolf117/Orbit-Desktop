package domain

import "time"

// Historial representa un registro en la tabla de historial de conversaciones
type Historial struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Nombre    string    `json:"nombre,omitempty"`
	Code      int64     `json:"code,omitempty"` // Código alfanumérico numérico único para enlazar con los JSONs locales
}

// HistorialRepository define los métodos para interactuar con la persistencia del historial
type HistorialRepository interface {
	GetHistorials() ([]Historial, error)
	CreateHistorial(h *Historial) error
	DeleteHistorial(id int) error
}
