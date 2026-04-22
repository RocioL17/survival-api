package repository

import (
    "encoding/json"
    "os"
    "sync"
    "survival-api/internal/models"
)

type JSONRepository struct {
    filePath string
    mu       sync.Mutex
}

func NewJSONRepository(path string) *JSONRepository {
    return &JSONRepository{filePath: path}
}

// Leer todos los casos
func (r *JSONRepository) GetCases() ([]models.Case, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    file, err := os.ReadFile(r.filePath)
    if err != nil {
        // Si el archivo no existe, devolvemos una lista vacía en lugar de error
        if os.IsNotExist(err) {
            return []models.Case{}, nil
        }
        return nil, err
    }

    var cases []models.Case
    if err := json.Unmarshal(file, &cases); err != nil {
        return nil, err
    }
    return cases, nil
}

// Guardar la lista completa (Reutilizable para cualquier cambio)
func (r *JSONRepository) WriteCases(cases []models.Case) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Convertimos la estructura a JSON con indentación para que sea legible
    data, err := json.MarshalIndent(cases, "", "  ")
    if err != nil {
        return err
    }

    // Escribimos el archivo (0644 son permisos estándar de lectura/escritura)
    return os.WriteFile(r.filePath, data, 0644)
}	