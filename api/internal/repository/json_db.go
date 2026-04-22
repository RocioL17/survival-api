package repository

import (
    "encoding/json"
    "os"
    "sync"
    "survival-api/models"
)

type JSONRepository struct {
    filePath string
    mu       sync.Mutex // Importante para evitar que dos personas escriban al mismo tiempo
}

func NewJSONRepository(path string) *JSONRepository {
    return &JSONRepository{filePath: path}
}

func (r *JSONRepository) GetAllUsers() ([]models.User, error) {
    r.mu.Lock()         // Bloqueamos para leer
    defer r.mu.Unlock() // Desbloqueamos al terminar

    file, _ := os.ReadFile(r.filePath)
    var users []models.User
    json.Unmarshal(file, &users)
    return users, nil
}