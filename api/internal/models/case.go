package models

import "encoding/json"

type Case struct {
	Name            string          `json:"name"`
	Age             int             `json:"age"`
	Gender          string          `json:"gender"`
	Latitud         float64         `json:"latitud"`
	Longitud        float64         `json:"longitud"`
	Zona            string          `json:"zone"`
	PuntosDeInteres []POI           `json:"interested_points"`
	Provincia       string          `json:"province"`
	Accidente       string          `json:"accident"`
	Choices         []string        `json:"choices"`
	ChoiceValue     []int           `json:"choice_value"`
	Historia        json.RawMessage `json:"story"`
}

type POI struct {
	Name       string
	Categories []string
}
