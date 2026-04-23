package handlers

import (
	"fmt"
    "encoding/json"
	"net/http"
	"survival-api/internal/services"
)

// type Data struct {
//  ID   string `json:"id"`
//  Name string `json:"name"`
// }

type CaseHandler struct {
	service *services.CaseService
}

func NewCaseHandler(service *services.CaseService) *CaseHandler {
	return &CaseHandler{service: service}
}

func (h *CaseHandler) MakeCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processedCase, err := h.service.CreateNewCase() 
    if err != nil {
        http.Error(w, "Error al generar el caso", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(processedCase)
}

func (h *CaseHandler) VerifyChoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Choice int `json:"choice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	result := h.service.RevisarPuntaje(body.Choice)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 - Page not found")
}
