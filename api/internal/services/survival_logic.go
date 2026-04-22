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

	edad, genero := GenerarDatosAleatorios()
	newCase.Age = edad
	newCase.Gender = genero
	log.Printf("[CreateNewCase] Datos aleatorios generados: edad=%d, genero=%s", edad, genero)

	log.Println("[CreateNewCase] Llamando a GenerarCase (maps)...")
	caseData := clients.GenerarCase()
	newCase.Latitud = caseData.Latitud
	newCase.Longitud = caseData.Longitud
	newCase.Zona = caseData.Zona
	newCase.Provincia = caseData.Provincia
	newCase.PuntosDeInteres = caseData.PuntosDeInteres
	log.Printf("[CreateNewCase] Ubicación obtenida: provincia=%s, zona=%s, lat=%v, lon=%v", newCase.Provincia, newCase.Zona, newCase.Latitud, newCase.Longitud)

	log.Printf("[CreateNewCase] Buscando accidente en dataset: provincia=%s, edad=%d, genero=%s", newCase.Provincia, newCase.Age, newCase.Gender)
	causa, err := s.BuscarAccidenteRandom(newCase.Provincia, newCase.Age, newCase.Gender)
	if err != nil {
		log.Printf("[CreateNewCase] Error buscando accidente: %v", err)
		return models.Case{}, err
	}
	newCase.Accidente = causa
	log.Printf("[CreateNewCase] Accidente encontrado: %s", causa)

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
	// Definimos las opciones de género
	generos := []string{"Masculino", "Femenino"}

	// rand.IntN(n) devuelve un número entre 0 y n-1
	edad := rand.IntN(73) + 18 // Genera de 0-72 y le suma 18 (Rango: 18 a 90)
	
	// Elegimos un índice al azar del arreglo de géneros
	generoElegido := generos[rand.IntN(len(generos))]

	return edad, generoElegido
}

func obtenerGrupoEdad(edad int) string {
    switch {
    case edad >= 15 && edad <= 19:
        return "15 a 19 años"
    case edad >= 20 && edad <= 24:
        return "20 a 24 años"
    case edad >= 25 && edad <= 29:
        return "25 a 29 años"
    case edad >= 30 && edad <= 34:
        return "30 a 34 años"
    case edad >= 35 && edad <= 39:
        return "35 a 39 años"
    default:
        return "Otros grupos" // Ajustar según los nombres reales de tu CSV
    }
}

// buscar en dataset para elegir un caso
func (s *CaseService) BuscarAccidenteRandom(provincia string, edad int, sexo string) (string, error) {
    f, err := os.Open("data/dataset.csv")
    if err != nil {
        return "", err
    }
    defer f.Close()

    reader := csv.NewReader(f)
    reader.Comma = ';' // IMPORTANTE: Tu CSV usa punto y coma
    records, err := reader.ReadAll()
    if err != nil {
        return "", err
    }

    // Convertimos la edad fija al formato del CSV
    grupoBuscado := obtenerGrupoEdad(edad)
    
    var causasPosibles []string

    for i, row := range records {
        if i == 0 { continue } // Saltamos la cabecera

        // Filtros basados en tus columnas:
        // row[2] = jurisdicion_residencia_nombre
        // row[6] = Sexo (o row[5] para sexo_id)
        // row[9] = grupo_edad
        
        matchProvincia := strings.TrimSpace(row[2]) == provincia
        matchSexo      := strings.TrimSpace(row[6]) == sexo
        matchEdad      := strings.TrimSpace(row[9]) == grupoBuscado

        if matchProvincia && matchSexo && matchEdad {
            // row[4] es cie10_clasificacion (el nombre de la causa)
            causasPosibles = append(causasPosibles, row[4])
        }
    }

	log.Printf("[BuscarAccidenteRandom] Filtros: provincia=%q, sexo=%q, grupoEdad=%q — resultados encontrados: %d", provincia, sexo, grupoBuscado, len(causasPosibles))
	if len(causasPosibles) == 0 {
		return "Sin datos para este perfil", nil
	}

	return causasPosibles[rand.IntN(len(causasPosibles))], nil
}

// armar prompt

// REVISAR PUNTAJE
func (s *CaseService) RevisarPuntaje(opcion int) (int) {
	cases, err := s.repo.GetCases()
    if err != nil {
        return 2
    }

	value := cases[0].ChoiceValue[opcion];
	return value;
}

// borrar datos

