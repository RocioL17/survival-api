package models

// User representa la estructura de un usuario en nuestro sistema.
// Los "tags" (entre comillas invertidas) le dicen a Go cómo
// convertir esto a JSON para tu interfaz.
type Case struct {
	Name            string   `json:"name"`
	Age             int      `json:"age"`
	Gender          string   `json:"gender"`
	Latitud         float64  `json:"latitud"`
	Longitud        float64  `json:"longitud"`
	Zona            string   `json:"zone"`
	PuntosDeInteres []POI    `json:"interested_points"`
	Provincia       string   `json:"province"`
	Accidente       string   `json:"accident"`
	Choices         []string `json:"choices"`
	ChoiceValue     []int    `json:"choice_value"`
	Historia        []byte   `json:"story"`
}

type POI struct {
	Name       string
	Categories []string
}
