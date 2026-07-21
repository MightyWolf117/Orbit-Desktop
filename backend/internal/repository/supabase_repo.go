package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"orbit-backend/internal/domain"
)

type SupabaseRepo struct {
	url      string
	password string // Se usará como apikey y Bearer token
	client   *http.Client
}

// Aseguramos en tiempo de compilación que cumple ambas interfaces
var _ domain.PersonalityRepository = (*SupabaseRepo)(nil)
var _ domain.HistorialRepository = (*SupabaseRepo)(nil)

func NewSupabaseRepo(url, password string) *SupabaseRepo {
	return &SupabaseRepo{
		url:      url,
		password: password,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *SupabaseRepo) GetPersonality() (*domain.Personality, error) {
	// Supabase REST endpoint para la tabla 'personalidad'.
	// Usamos query params para obtener el registro activo si lo hubiera, 
	// o por defecto traemos el primer registro (limit=1).
	reqUrl := fmt.Sprintf("%s?select=*&limit=1", r.url)
	
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[Supabase GET GetPersonality] URL: %s", reqUrl)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}
	
	log.Printf("[Supabase GET GetPersonality] Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var personalities []domain.Personality
	if err := json.Unmarshal(bodyBytes, &personalities); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de supabase: %w", err)
	}

	if len(personalities) == 0 {
		return nil, fmt.Errorf("no se encontraron personalidades en la base de datos")
	}

	return &personalities[0], nil
}

func (r *SupabaseRepo) GetPersonalityByID(id int) (*domain.Personality, error) {
	reqUrl := fmt.Sprintf("%s?id=eq.%d&select=*&limit=1", r.url, id)
	
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var personalities []domain.Personality
	if err := json.Unmarshal(bodyBytes, &personalities); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de supabase: %w", err)
	}

	if len(personalities) == 0 {
		return nil, fmt.Errorf("no se encontró la personalidad con ID %d", id)
	}

	return &personalities[0], nil
}

func (r *SupabaseRepo) GetAllPersonalities() ([]domain.Personality, error) {
	reqUrl := fmt.Sprintf("%s?select=*", r.url)
	
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[Supabase GET GetAllPersonalities] URL: %s", reqUrl)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}
	
	log.Printf("[Supabase GET GetAllPersonalities] Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var personalities []domain.Personality
	if err := json.Unmarshal(bodyBytes, &personalities); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de supabase: %w", err)
	}

	return personalities, nil
}

func (r *SupabaseRepo) CreatePersonality(p *domain.Personality) error {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error codificando personalidad: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, r.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")
	// Prefer: return=representation devuelve el registro insertado, 
	// pero si solo queremos crear podemos omitirlo. Lo ponemos por buena práctica.
	req.Header.Set("Prefer", "return=representation")
	
	log.Printf("[Supabase POST CreatePersonality] URL: %s, Payload: %s", r.url, string(jsonData))

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}
	
	log.Printf("[Supabase POST CreatePersonality] Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (r *SupabaseRepo) UpdatePersonality(id int, p *domain.Personality) error {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error codificando personalidad: %w", err)
	}

	reqUrl := fmt.Sprintf("%s?id=eq.%d", r.url, id)
	req, err := http.NewRequest(http.MethodPatch, reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation") // Para que retorne la data modificada

	log.Printf("[Supabase PATCH UpdatePersonality] URL: %s, Payload: %s", reqUrl, string(jsonData))

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}
	
	log.Printf("[Supabase PATCH UpdatePersonality] Response Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Historial Methods

func (r *SupabaseRepo) GetHistorials() ([]domain.Historial, error) {
	reqUrl := fmt.Sprintf("%s/historial?select=*&order=created_at.desc", r.url)
	
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var historiales []domain.Historial
	if err := json.Unmarshal(bodyBytes, &historiales); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return historiales, nil
}

func (r *SupabaseRepo) CreateHistorial(h *domain.Historial) error {
	jsonData, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("error codificando historial: %w", err)
	}

	reqUrl := fmt.Sprintf("%s/historial", r.url)
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Actualizar el ID generado
	var inserted []domain.Historial
	if err := json.Unmarshal(bodyBytes, &inserted); err == nil && len(inserted) > 0 {
		h.ID = inserted[0].ID
		h.CreatedAt = inserted[0].CreatedAt
	}

	return nil
}

func (r *SupabaseRepo) DeleteHistorial(id int) error {
	reqUrl := fmt.Sprintf("%s/historial?id=eq.%d", r.url, id)
	
	req, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("apikey", r.password)
	req.Header.Set("Authorization", "Bearer "+r.password)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error al hacer petición a supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase devolvió estado no exitoso: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
