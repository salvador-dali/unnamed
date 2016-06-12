// Package structs contains structs for all models
package brands

import (
	"encoding/json"
	"github.com/salvador-dali/unnamed/structs"
	"log"
	"net/http"
	"strconv"
	"database/sql"
)

// getIntegerID checks whether the string representation of an ID is positive integer
// If not, returns 0, sends back an empty json with 404 status code
func getIntegerID(w http.ResponseWriter, idString string) int {
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return 0
	}
	return id
}


func GetAllBrands(db *sql.DB) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
    return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
        w.Header().Set("Content-Type", "application/javascript")

		brands := []*structs.Brand{}
		rows, err := db.Query("SELECT id, name FROM brands")
		defer rows.Close()
		for rows.Next() {
			brand := structs.Brand{}
			if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
				log.Fatal(err)
			}
			brands = append(brands, &brand)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		json, _ := json.Marshal(brands)
		w.Write(json)
    }
}

func GetBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := getIntegerID(w, ps["id"])
	if id <= 0 {
		return
	}

	json, _ := json.Marshal(structs.Id{id})
	w.Write(json)
}

func CreateBrand(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
