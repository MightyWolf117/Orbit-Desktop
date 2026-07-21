package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"orbit-backend/internal/domain"
)

type aiService struct {
	apiKey   string
	persRepo domain.PersonalityRepository
	histRepo domain.HistorialRepository
}

func NewAIService(apiKey string, persRepo domain.PersonalityRepository, histRepo domain.HistorialRepository) domain.AIService {
	return &aiService{
		apiKey:   apiKey,
		persRepo: persRepo,
		histRepo: histRepo,
	}
}

func (s *aiService) GenerateResponse(ctx context.Context, req domain.ChatRequest) (*domain.ChatResponse, error) {
	var personality *domain.Personality
	var err error

	if req.PersonalityID != nil {
		personality, err = s.persRepo.GetPersonalityByID(*req.PersonalityID)
		if err != nil {
			log.Printf("Error obteniendo personalidad por ID %d: %v", *req.PersonalityID, err)
		}
	}

	// Fallback si no hay ID, o si falló la búsqueda por ID
	if personality == nil {
		personality, err = s.persRepo.GetPersonality()
		if err != nil {
			log.Printf("Error obteniendo personalidad por defecto: %v", err)
			personality = &domain.Personality{
				Prompt: "Eres un asistente virtual útil e inteligente.",
			}
		}
	}

	// Inicializar el cliente de Google Gemini
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		return nil, fmt.Errorf("error inicializando cliente gemini: %w", err)
	}
	defer client.Close()

	modelName := "gemini-flash-latest"
	if req.Model != "" {
		modelName = req.Model
	}
	model := client.GenerativeModel(modelName)

	if req.Temperature != nil {
		model.Temperature = req.Temperature
	}

	// Configuramos las instrucciones del sistema (la personalidad)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(personality.Prompt)},
	}

	// Iniciar una sesión de chat
	cs := model.StartChat()

	// Añadir el historial de mensajes a la sesión
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("la solicitud de chat no contiene mensajes")
	}

	// Agregar el historial (excluyendo el último mensaje que es el prompt actual)
	for i := 0; i < len(req.Messages)-1; i++ {
		msg := req.Messages[i]
		var role string
		if msg.Role == "model" || msg.Role == "assistant" {
			role = "model"
		} else {
			role = "user"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// El último mensaje del usuario
	lastMsg := req.Messages[len(req.Messages)-1]

	// Enviar mensaje a Gemini
	resp, err := cs.SendMessage(ctx, genai.Text(lastMsg.Content))
	if err != nil {
		return nil, fmt.Errorf("error generando respuesta de gemini: %w", err)
	}

	// Extraer el texto de la respuesta
	var responseText string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			responseText = string(text)
		}
	}

	// Generar título si fue solicitado y tenemos mensajes
	var title string
	var historial *domain.Historial
	if req.GenerateTitle && len(req.Messages) > 0 {
		firstMsg := req.Messages[0].Content
		titlePrompt := fmt.Sprintf("Genera un título corto (máximo 4 palabras) para una conversación que comienza con este mensaje: \"%s\". Responde ÚNICAMENTE con el título, sin comillas ni texto adicional.", firstMsg)

		titleModel := client.GenerativeModel("gemini-flash-latest")
		titleResp, err := titleModel.GenerateContent(ctx, genai.Text(titlePrompt))
		if err == nil && len(titleResp.Candidates) > 0 && len(titleResp.Candidates[0].Content.Parts) > 0 {
			if t, ok := titleResp.Candidates[0].Content.Parts[0].(genai.Text); ok {
				title = string(t)
				// Guardar en Supabase si tenemos código de chat
				if req.ChatCode != nil {
					newHistorial := &domain.Historial{
						Nombre: title,
						Code:   *req.ChatCode,
					}
					if err := s.histRepo.CreateHistorial(newHistorial); err != nil {
						log.Printf("Error creando historial en base de datos: %v", err)
					} else {
						historial = newHistorial
					}
				}
			}
		} else {
			log.Printf("Error o respuesta vacía generando título: %v", err)
		}
	}

	return &domain.ChatResponse{
		Response:  responseText,
		Title:     title,
		Historial: historial,
	}, nil
}
