package clients

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"survival-api/internal/models"
	"time"
)

// ── Estructuras del REQUEST a Gemini ─────────────────────────────────────────

// se define un GroqRequest que es el formato que espera la API de Groq para generar una respuesta basada en un prompt del sistema y un mensaje del usuario.
// Incluye el modelo a usar, la temperatura para la generación, y el formato de respuesta esperado (en este caso, un objeto JSON).
type GroqRequest struct {
	Model          string         `json:"model"`
	Messages       []Message      `json:"messages"`        // El prompt del sistema y el mensaje del usuario
	Temperature    float64        `json:"temperature"`     // Controla la creatividad de la respuesta, valores más altos generan respuestas más creativas
	ResponseFormat ResponseFormat `json:"response_format"` // Especifica que queremos la respuesta en formato JSON
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

// ── Estructuras de la RESPUESTA de Groq ────────────────────────────────────

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ── Estructura del JSON que pedimos DENTRO de la respuesta ───────────────────

type Opcion struct {
	ID          string `json:"id"`          // Identificador de la opción (A, B, C)
	Texto       string `json:"texto"`       // Texto descriptivo de la opción
	EsSalvacion bool   `json:"esSalvacion"` // Indica si esta opción es la correcta para salvar al jugador
}

type HistoriaResponse struct {
	Historia string   `json:"historia"` // Narración de la historia en 3-4 oraciones
	Pregunta string   `json:"pregunta"` // Pregunta de decisión para el jugador
	Nombre   string   `json:"nombre"`
	Opciones []Opcion `json:"opciones"` // Lista de opciones de decisión, exactamente 3, con una sola opción de salvación
}

// ── Perfil del jugador ESTO ME PASA RO

type PerfilJugador struct {
	Edad        int    `json:"edad"`
	Sexo        string `json:"sexo"`
	Ubicacion   string `json:"ubicacion"`
	CausaMuerte string `json:"causaMuerte"`
	Transito    string `json:"transito"`
}

// ── Configuración de Groq  ────────────────────────────────────────────────────────────

const (
	groqURL          = "https://api.groq.com/openai/v1/chat/completions"
	defaultGroqModel = "llama-3.1-8b-instant"
	maxRetries       = 2
)

var systemPrompt = `Sos el narrador de un juego de supervivencia llamado "Vive o Muere".
Tu tarea es generar una historia corta y dramática basada en el perfil y situación del jugador.

REGLAS:
- Respondé ÚNICAMENTE con un JSON válido, sin texto extra ni bloques de código.
- El JSON debe tener exactamente esta estructura:
{
  "nombre": "Nombre realista según edad, sexo y país de la ubicación",
  "historia": "Narración en 3-4 oraciones en segunda persona (vos/tú). Mencioná el lugar específico.",
  "pregunta": "Pregunta de decisión urgente y concreta.",
  "opciones": [
    { "id": "A", "texto": "...", "esSalvacion": false },
    { "id": "B", "texto": "...", "esSalvacion": true },
    { "id": "C", "texto": "...", "esSalvacion": false }
  ]
}
- Exactamente 3 opciones, exactamente 1 con esSalvacion=true.
- La opción correcta debe variar de posición (no siempre la B).
- El tono debe ser tenso y realista, como un documental de supervivencia.`

func buildUserPrompt(p models.Case) string {
	// Contexto de zona
	zonaCtx := "zona urbana con acceso a servicios"
	if p.Zona == "rural" {
		zonaCtx = "zona rural alejada, con difícil acceso a servicios de emergencia"
	}

	// Puntos de interés cercanos
	puntosCtx := "No hay puntos de interés relevantes cercanos."
	nombres := make([]string, len(p.PuntosDeInteres))
	for i, poi := range p.PuntosDeInteres {
		if len(poi.Categories) > 0 {
			nombres[i] = poi.Name + " (" + strings.Join(poi.Categories, "/") + ")"
		} else {
			nombres[i] = poi.Name
		}
	}
	puntosCtx = "Lugares cercanos que el jugador podría usar: " + strings.Join(nombres, ", ") + "."

	return fmt.Sprintf(`Generá una historia para este jugador:

PERFIL:
- Edad: %d años
- Sexo: %s
- Provincia/Región: %s
- Coordenadas: lat %s, lon %s
- Tipo de zona: %s
- Accidente/situación: %s

ENTORNO:
- %s

INSTRUCCIONES:
- Usá la provincia y el tipo de zona para darle autenticidad al lugar.
- Si es zona rural, las opciones deben reflejar la lejanía (ej: caminar hasta el pueblo, llamar por radio, usar recursos naturales).
- Si es zona urbana, las opciones pueden incluir los puntos de interés cercanos (hospitales, iglesias, farmacias, etc.).
- La opción correcta debe ser la más lógica dado el contexto, no necesariamente la más obvia.
- Generá un nombre realista para el personaje según su edad, sexo y región.`,
		p.Age,
		p.Gender,
		p.Provincia,
		p.Latitud,
		p.Longitud,
		zonaCtx,
		p.Accidente,
		puntosCtx,
	)
}

// esto es para manejar variables de entorno desde un archivo .env
func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// esto es para manejar reintentos en caso de errores temporales al llamar a la API de Groq,
// como timeouts o errores 500. Se intentará hasta maxRetries veces antes de devolver un error definitivo.
func LlamarGroq(perfil models.Case) (*HistoriaResponse, error) {
	if err := loadDotEnv(".env"); err != nil {
		return nil, fmt.Errorf("error cargando .env: %w", err)
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GROQ_API_KEY no está definida")
	}

	model := strings.TrimSpace(os.Getenv("GROQ_MODEL"))
	if model == "" {
		model = defaultGroqModel
	}

	reqBody := GroqRequest{
		Model:          model,
		Temperature:    0.9,
		ResponseFormat: ResponseFormat{Type: "json_object"},
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: buildUserPrompt(perfil)},
		},
	}

	var lastErr error
	for intento := 1; intento <= maxRetries; intento++ {
		historia, err := doRequest(groqURL, apiKey, reqBody)
		if err == nil {
			return historia, nil
		}
		lastErr = err
		fmt.Printf("[Groq] intento %d fallido: %v\n", intento, err)
		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("groq falló tras %d intentos: %w", maxRetries, lastErr)
}

// esta función hace la llamada HTTP a la API de Groq, maneja la respuesta,
// extrae el JSON de la respuesta (incluso si viene con texto adicional o formateado como código),
// y valida que el JSON tenga la estructura esperada antes de devolverlo.
func doRequest(url string, apiKey string, reqBody GroqRequest) (*HistoriaResponse, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error serializando request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creando request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en HTTP POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq devolvió status %d: %s", resp.StatusCode, string(body))
	}

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de Groq: %w", err)
	}
	if len(groqResp.Choices) == 0 {
		return nil, errors.New("groq devolvió una respuesta vacía")
	}

	textoJSON, err := extractJSON(groqResp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("respuesta sin JSON parseable: %w", err)
	}

	var historia HistoriaResponse
	if err := json.Unmarshal([]byte(textoJSON), &historia); err != nil {
		return nil, fmt.Errorf("el JSON de Groq no es válido: %w", err)
	}

	if err := validarHistoria(&historia); err != nil {
		return nil, err
	}

	return &historia, nil
}

// esta función intenta extraer un objeto JSON válido de un texto que puede contener texto adicional, formateo de código, o estar mal formateado.
// Primero intenta validar el texto completo, luego busca bloques de código, y finalmente intenta extraer un substring entre llaves.
func extractJSON(text string) (string, error) {
	trimmed := strings.TrimSpace(text)

	if json.Valid([]byte(trimmed)) {
		return trimmed, nil
	}

	if strings.Contains(trimmed, "```") {
		parts := strings.Split(trimmed, "```")
		for _, part := range parts {
			candidate := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(part), "json"))
			if candidate != "" && json.Valid([]byte(candidate)) {
				return candidate, nil
			}
		}
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		candidate := strings.TrimSpace(trimmed[start : end+1])
		if json.Valid([]byte(candidate)) {
			return candidate, nil
		}
	}

	return "", errors.New("no se encontró un objeto JSON válido en la respuesta")
}

// esta función valida que la historia generada por Groq tenga la estructura correcta:
// - historia no vacía
// - pregunta no vacía
// - exactamente 3 opciones
// - exactamente 1 opción de salvación
func validarHistoria(h *HistoriaResponse) error {
	if h.Historia == "" {
		return errors.New("campo 'historia' vacío")
	}
	if h.Pregunta == "" {
		return errors.New("campo 'pregunta' vacío")
	}
	if len(h.Opciones) != 3 {
		return fmt.Errorf("se esperaban 3 opciones, llegaron %d", len(h.Opciones))
	}

	salvaciones := 0
	for _, op := range h.Opciones {
		if op.EsSalvacion {
			salvaciones++
		}
	}
	if salvaciones != 1 {
		return fmt.Errorf("debe haber exactamente 1 opción de salvación, hay %d", salvaciones)
	}

	return nil
}

// esto se tiene que remplazar con toda la info que se obtenga del dataset y la api de mapas

func comoUsar() {
	perfil := models.Case{
		Age:       22,
		Gender:    "Femenino",
		Zona:      "Misiones, Argentina",
		Accidente: "Accidente en la ruta",
		Latitud:   -34.3917,
		Longitud:  -58.8731,
		Provincia: "Misiones",
	}

	historia, err := LlamarGroq(perfil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(historia, "", "  ")
	if err != nil {
		fmt.Println("Error serializando salida JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
