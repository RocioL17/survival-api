package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"survival-api/internal/clients"
	"survival-api/internal/models"
	"survival-api/internal/repository"
)

type CaseService struct {
	repo *repository.JSONRepository
}

func NewCaseService(repo *repository.JSONRepository) *CaseService {
	return &CaseService{repo: repo}
}

// CREACIÓN DE UN CASO
// En internal/services/case_service.go
func (s *CaseService) CreateNewCase() (models.Case, error) {
	log.Println("[CreateNewCase] Iniciando creación de caso")
	var newCase models.Case

	log.Println("[CreateNewCase] Llamando a GenerarCase (maps)...")
	caseData := clients.GenerarCase()
	newCase.Latitud = caseData.Latitud
	newCase.Longitud = caseData.Longitud
	newCase.Zona = caseData.Zona
	newCase.Provincia = caseData.Provincia
	newCase.PuntosDeInteres = caseData.PuntosDeInteres
	log.Printf("[CreateNewCase] Ubicación obtenida: provincia=%s, zona=%s, lat=%v, lon=%v", newCase.Provincia, newCase.Zona, newCase.Latitud, newCase.Longitud)

	causa, edad, genero, err := s.BuscarAccidenteRandom(newCase.Provincia)
	if err != nil {
		log.Printf("[CreateNewCase] Error buscando accidente: %v", err)
		return models.Case{}, err
	}
	newCase.Age = edad
	newCase.Gender = genero
	newCase.Accidente = causa
	log.Printf("[CreateNewCase] Accidente encontrado: %s (edad=%d, genero=%s)", causa, edad, genero)

	perfil := models.Case{
		Age:             newCase.Age,
		Gender:          newCase.Gender,
		Zona:            newCase.Zona,
		PuntosDeInteres: newCase.PuntosDeInteres,
		Provincia:       newCase.Provincia,
		Accidente:       newCase.Accidente,
	}

	log.Println("[CreateNewCase] Llamando a LlamarGroq (IA)...")
	historia, err := clients.LlamarGroq(perfil)
	if err != nil {
		log.Printf("[CreateNewCase] Error en LlamarGroq: %v", err)
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	log.Println("[CreateNewCase] Historia generada por IA correctamente")

	choices := make([]string, len(historia.Opciones))
	choiceValues := make([]int, len(historia.Opciones))
	for i, op := range historia.Opciones {
		choices[i] = op.Texto
		if op.EsSalvacion {
			choiceValues[i] = 1
		} else {
			choiceValues[i] = 0
		}
	}
	newCase.Choices = choices
	newCase.ChoiceValue = choiceValues

	newCase.Name = historia.Nombre
	log.Printf("[CreateNewCase] Opciones procesadas: %d opciones", len(choices))

	historiaResto := struct {
		Historia string `json:"historia"`
		Pregunta string `json:"pregunta"`
	}{
		Historia: historia.Historia,
		Pregunta: historia.Pregunta,
	}
	historiaBytes, err := json.Marshal(historiaResto)
	if err != nil {
		log.Printf("[CreateNewCase] Error serializando historia: %v", err)
		return models.Case{}, fmt.Errorf("error serializando historia: %w", err)
	}
	newCase.Historia = historiaBytes

	log.Println("[CreateNewCase] Guardando caso en repositorio...")
	cases, err := s.repo.GetCases()
	if err != nil {
		log.Printf("[CreateNewCase] No se encontraron casos previos, inicializando lista vacía: %v", err)
		s.repo.WriteCases([]models.Case{})
	}

	cases = append(cases, newCase)
	err = s.repo.WriteCases(cases)
	if err != nil {
		log.Printf("[CreateNewCase] Error escribiendo casos: %v", err)
		return models.Case{}, err
	}

	log.Println("[CreateNewCase] Caso creado y guardado exitosamente")
	return newCase, nil
}

// random edad sexo
func GenerarDatosAleatorios() (int, string) {
	generos := []string{"masculino", "femenino"}
	edad := rand.IntN(73) + 18 // Rango: 18 a 90
	generoElegido := generos[rand.IntN(len(generos))]
	return edad, generoElegido
}

func obtenerGrupoEdad(edad int) string {
	switch {
	case edad <= 14:
		return "01.De a 0  a 14 anios"
	case edad <= 34:
		return "02.De 15 a 34 anios"
	case edad <= 54:
		return "03.De 35 a 54 anios"
	case edad <= 74:
		return "04.De 55 a 74 anios"
	default:
		return "05.De 75 anios y mas"
	}
}

// buscar en dataset para elegir un caso, generando edad/género aleatoriamente y reintentando hasta encontrar resultados
func (s *CaseService) BuscarAccidenteRandom(provincia string) (causa string, edad int, genero string, err error) {
	f, err := os.Open("data/dataset.csv")
	if err != nil {
		return "", 0, "", err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return "", 0, "", err
	}

	for intento := 1; intento <= 2000; intento++ {
		edad, genero = GenerarDatosAleatorios()
		grupoBuscado := obtenerGrupoEdad(edad)

		var causasPosibles []string
		for i, row := range records {
			if i == 0 {
				continue
			}
			if strings.TrimSpace(row[2]) == provincia &&
				strings.TrimSpace(row[6]) == genero &&
				strings.TrimSpace(row[9]) == grupoBuscado {
				causasPosibles = append(causasPosibles, row[4])
			}
		}

		log.Printf("[BuscarAccidenteRandom] Intento %d — provincia=%q, genero=%q, grupoEdad=%q — resultados: %d", intento, provincia, genero, grupoBuscado, len(causasPosibles))

		if len(causasPosibles) > 0 {
			return causasPosibles[rand.IntN(len(causasPosibles))], edad, genero, nil
		}
	}

	return "", 0, "", fmt.Errorf("no se encontraron accidentes para provincia=%q después de 20 intentos", provincia)
}

// armar prompt

// REVISAR PUNTAJE
func (s *CaseService) RevisarPuntaje(opcion int) int {
	cases, err := s.repo.GetCases()
	if err != nil {
		return 2
	}

	value := cases[0].ChoiceValue[opcion]
	return value
}

// borrar datos
