// Package structs contains structs for all models
package brands

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/salvador-dali/unnamed/errorCodes"
	"github.com/salvador-dali/unnamed/helpers"
	"github.com/salvador-dali/unnamed/structs"
	"log"
	"net/http"
)

const maximumNameLength = 40

// GetAllBrands returns the (id, name)
func GetAllBrands(db *sql.DB) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "application/javascript")

		brands := []*structs.Brand{}
		rows, err := db.Query("SELECT id, name FROM brands")
		if err != nil {
			log.Fatal(err)
		}
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

		json, err := json.Marshal(brands)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(json)
	}
}

// GetBrand returns full information about a brand
func GetBrand(db *sql.DB) func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		w.Header().Set("Content-Type", "application/javascript")

		id, brand := helpers.ValidateId(w, ps["id"]), structs.Brand{}
		if id == nil {
			return
		}

		if err := db.QueryRow("SELECT id, name, issued_at FROM brands WHERE id = $1", id).Scan(&brand.Id, &brand.Name, &brand.Issued_at); err != nil {
			if err == sql.ErrNoRows {
				json, err := json.Marshal(structs.ErrorCode{errorCodes.IdNotExist})
				if err != nil {
					log.Fatal(err)
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(json)
				return
			} else {
				log.Fatal(err)
			}
		}

		json, err := json.Marshal(brand)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(json)
		return
	}
}

// CreateBrand creates a brand with a specific name
func CreateBrand(db *sql.DB) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "application/javascript")

		if !helpers.isValidateFormLength(w, r, 1) {
			return
		}

		name := helpers.ValidateName(r.PostFormValue("name"), maximumNameLength)
		if name == nil {
			return
		}

		elementId := 0
		err := db.QueryRow("INSERT INTO brands (name) VALUES($1) RETURNING id", name).Scan(&elementId)
		if err != nil {
			if errPg, ok := err.(*pq.Error); ok && string(errPg.Code) == "23505" {
				// 23505 is : duplicate key value violates unique constraint. Names can't be the same
				json, err := json.Marshal(structs.ErrorCode{errorCodes.DuplicateName})
				if err != nil {
					log.Fatal(err)
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write(json)
				return
			}
			log.Fatal(err)
		}

		json, err := json.Marshal(structs.Id{int(elementId)})
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(json)
	}
}

func UpdateBrand(db *sql.DB) func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		w.Header().Set("Content-Type", "application/javascript")

		id, brand := getIntegerID(w, ps["id"]), structs.Brand{}
		if id <= 0 {
			return
		}

	}
}
