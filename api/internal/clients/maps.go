package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	"log"
	
)



func getVariable() string{
	err:= godotenv.Load("../../.env")

	if err!=nil {
		log.Fatal("Error intentando cargar el .env")
		log.Fatal(err)

	}

	variable := os.Getenv("TT_API_KEY")

	// fmt.Println(variable)

	return variable

}



var apiKey = getVariable()


// ─── Coordenada random ─────────────────────────────
func randomCoord() (float64, float64) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for {
        lat := -55.05 + r.Float64()*33.27
        lon := -73.57 + r.Float64()*19.93

        // --- FILTROS MANUALES RÁPIDOS ---
        
        // 1. Evitar el rincón del Atlántico (si es muy al sur y muy al este, es mar)
        if lat < -42.0 && lon > -63.0 {
            continue 
        }

        // 2. Evitar Chile profundo (si es muy al oeste en latitudes centrales)
        if lat < -30.0 && lat > -45.0 && lon < -71.5 {
            continue
        }

        // 3. Evitar Uruguay/Brasil (si es latitud de Buenos Aires/Litoral pero muy al este)
        if lat > -35.0 && lat < -30.0 && lon > -58.0 {
            continue
        }

        return lat, lon
    }
}
func randomCoordArgentina() (float64, float64, string) {
	for {
		lat, lon := randomCoord()

		// snap opcional (mejora precisión)
		lat, lon = snapToRoad(lat, lon)

		provincia, country := getProvincia(lat, lon)

		// Validación fuerte
		if country == "AR" && provincia != "" && provincia != "desconocida" {
			return lat, lon, provincia
		}

		fmt.Println("Reintentando coordenada...")
	}
}

// ─── Snap to Road ──────────────────────────────────
func snapToRoad(lat, lon float64) (float64, float64) {
	url := fmt.Sprintf(
		"https://api.tomtom.com/routing/1/snapToRoads/%f,%f.json?key=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return lat, lon
	}
	defer resp.Body.Close()

	var data struct {
		Points []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"points"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return lat, lon
	}

	if len(data.Points) > 0 {
		return data.Points[0].Latitude, data.Points[0].Longitude
	}

	return lat, lon
}

// ─── POIs ──────────────────────────────────────────
type POI struct {
	Name       string
	Categories []string
}

func getPOIs(lat, lon float64) ([]POI, error) {
	url := fmt.Sprintf(
		"https://api.tomtom.com/search/2/nearbySearch/.json?lat=%f&lon=%f&radius=10000&limit=8&key=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			POI struct {
				Name       string   `json:"name"`
				Categories []string `json:"categories"`
			} `json:"poi"`
		} `json:"results"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	var pois []POI
	for _, r := range data.Results {
		pois = append(pois, POI{
			Name:       r.POI.Name,
			Categories: r.POI.Categories,
		})
	}

	return pois, nil
}

// ─── Provincia ─────────────────────────────────────
// func getProvincia(lat, lon float64) string {
// 	url := fmt.Sprintf(
// 		"https://api.tomtom.com/search/2/reverseGeocode/%f,%f.json?key=%s",
// 		lat, lon, apiKey,
// 	)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "desconocida"
// 	}
// 	defer resp.Body.Close()

// 	var data struct {
// 		Addresses []struct {
// 			Address struct {
// 				CountrySubdivision string `json:"countrySubdivision"`
// 			} `json:"address"`
// 		} `json:"addresses"`
// 	}

// 	json.NewDecoder(resp.Body).Decode(&data)

// 	if len(data.Addresses) > 0 {
// 		return data.Addresses[0].Address.CountrySubdivision
// 	}

// 	return "desconocida"
// }

func getProvincia(lat, lon float64) (string, string) {
	url := fmt.Sprintf(
		"https://api.tomtom.com/search/2/reverseGeocode/%f,%f.json?key=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "desconocida", ""
	}
	defer resp.Body.Close()

	var data struct {
		Addresses []struct {
			Address struct {
				CountryCode        string `json:"countryCode"`
				CountrySubdivision string `json:"countrySubdivision"`
			} `json:"address"`
		} `json:"addresses"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	if len(data.Addresses) > 0 {
		addr := data.Addresses[0].Address
		return addr.CountrySubdivision, addr.CountryCode
	}

	return "desconocida", ""
}

func getZona(pois []POI) string {
	if len(pois) > 3 {
		return "urbano"
	}
	return "rural"
}

func printPOIs(pois []POI) {
	fmt.Println("POIs cercanos:")

	if len(pois) == 0 {
		fmt.Println("  - No se encontraron")
		return
	}

	for _, p := range pois {
		fmt.Printf("  - %s (%v)\n", p.Name, p.Categories)
	}
}

// ─── MAIN ──────────────────────────────────────────
func main() {
	// lat, lon := randomCoord()
	lat, lon, provincia := randomCoordArgentina()

	
	// lat, lon = snapToRoad(lat, lon)

	pois, _ := getPOIs(lat, lon)
	// provincia := getProvincia(lat, lon)

	
	// Clasificación simple
	
	zona := getZona(pois)


	fmt.Println("Coordenadas:", lat, lon)
	fmt.Println("Provincia:", provincia)
	fmt.Println("Zona:", zona)

	printPOIs(pois)

    

}