package services

import (
	"math/rand/v2"
	"survival-api/internal/models"
	"survival-api/internal/repository"
	"encoding/csv"
	"os"
	"strings"
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
	var newCase models.Case	
    // 1. Generamos los datos aleatorios sobre el objeto que recibimos
    edad, genero := GenerarDatosAleatorios()
    newCase.Age = edad
    newCase.Gender = genero
    
    // latitu longitu

	// llamada a maps

	// buscar en dataset función
	causa, err := s.BuscarAccidenteRandom(newCase.Provincia, newCase.Age, newCase.Gender)
    if err != nil {
        return models.Case{}, err
    }
    
    newCase.Accidente = causa

	// llamar a ia
	perfil := {
		Edad:        newCase.Age,
		Sexo:        newCase.Gender,
		Ubicacion:   newCase.Provincia,
		CausaMuerte: "accidente en la ruta",
		Transito:    "Muchos autos, clima lluvioso, visibilidad reducida",
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
    
	// Limpieza del json y escritura de los datos nuevos
	cases, err := s.repo.GetCases()
    if err != nil {
        s.repo.WriteCases([]models.Case{}) //si hay algún caso previo se borra
    }

    cases = append(cases, newCase)
    err = s.repo.WriteCases(cases)
    if err != nil {
        return models.Case{}, err
    }

    return newCase, nil
}


// random edad sexo
func GenerarDatosAleatorios() (int, string) {
	// Definimos las opciones de género
	generos := []string{"Masculino", "Femenino", "No binario", "Prefiero no decirlo"}

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

    if len(causasPosibles) == 0 {
        return "Sin datos para este perfil", nil
    }

    // Retornamos una causa aleatoria de las que cumplen el filtro
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

