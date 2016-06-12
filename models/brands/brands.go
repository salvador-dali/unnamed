// Package structs contains structs for all models
package brands

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/salvador-dali/unnamed/structs"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func validateName(name string, n int) (string, error) {
	if name == "" {
		return "", errors.New("Wrong string provided")
	}

	name = strings.TrimSpace(name)
	if len(name) > 0 && len(name) <= n {
		return name, nil
	}

	return "", errors.New("Wrong string provided")
}

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

		id, brand := getIntegerID(w, ps["id"]), structs.Brand{}
		if id <= 0 {
			return
		}

		db.QueryRow("SELECT id, name, issued_at FROM brands WHERE id = $1", id).Scan(&brand)

		json, err := json.Marshal(brand)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id)
		fmt.Println(brand)
		fmt.Println(json)

		w.Write(json)
	}
}

func CreateBrand(db *sql.DB) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "application/javascript")

		r.ParseForm()
		if len(r.Form) != 1 {
			// TODO provide information why failed
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO remove constant
		name, err := validateName(r.PostFormValue("name"), 40)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		elementId := 0
		err = db.QueryRow("INSERT INTO brands (name) VALUES($1) RETURNING id", name).Scan(&lastInsertId)
		if err != nil {
			if errPg, ok := err.(*pq.Error); ok && string(errPg.Code) == "23505" {
				// 23505 is a code for: duplicate key value violates unique constraint.
				// Names can't be the same
				json, err := json.Marshal(structs.ErrorCode{100})
				if err != nil {
					log.Fatal(err)
				}
				w.Write(json)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Fatal(err)
		}

		json, err := json.Marshal(structs.Id{int(elementId)})
		if err != nil {
			log.Fatal(err)
		}
		w.Write(json)
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
