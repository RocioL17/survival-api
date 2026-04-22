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
    service *services.CaseService // El receptor tiene guardado su motor (el servicio)
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

func VerifyChoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 - Page not found")
}
