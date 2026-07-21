package domain

// Personality representa la personalidad del asistente recuperada de Supabase
type Personality struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"nombre,omitempty"`
	Description string  `json:"descripcion_corta,omitempty"`
	Prompt      string  `json:"instrucciones,omitempty"` // Instrucciones base (ej. "Eres un asistente amable...")
	Image       *string `json:"image,omitempty"`
}

// PersonalityRepository define los métodos para interactuar con la persistencia de personalidades
type PersonalityRepository interface {
	GetPersonality() (*Personality, error)
	GetPersonalityByID(id int) (*Personality, error)
	GetAllPersonalities() ([]Personality, error)
	CreatePersonality(p *Personality) error
	UpdatePersonality(id int, p *Personality) error
}
