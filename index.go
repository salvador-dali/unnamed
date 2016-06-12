package main

import (
	"encoding/json"
	"fmt"
	"database/sql"
	"github.com/dimfeld/httptreemux"
	"github.com/salvador-dali/unnamed/config"
	"github.com/salvador-dali/unnamed/structs"
	"log"
	"net/http"
	"strconv"
	_ "github.com/lib/pq"
)

// Cfg stores information about all environment variables
var Cfg config.Config

// Db is a connection to a PSQL database
var Db *sql.DB

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

//----------- Brands
func GetAllBrands(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	brands := []*structs.Brand{}
	rows, err := Db.Query("SELECT id, name FROM brands")
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

func GetBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := getIntegerID(w, ps["id"])
	if id <= 0 {
		return
	}

	json, _ := json.Marshal(Id{id})
	w.Write(json)
}

func CreateBrand(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	// Create a router, initialize db connection and config
	router, Cfg := httptreemux.New(), config.GetConfig()
	api, dbURL := router.NewGroup("/api/v1"), fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", Cfg.DbUser, Cfg.DbPass, Cfg.DbHost, Cfg.DbPort, Cfg.DbName)
	if db, err := sql.Open("postgres", dbURL); err != nil {
		log.Fatal(err)
	} else {
		Db = db
	}

	// Brands
	api.GET("/brands", GetAllBrands)
	api.GET("/brands/:id", GetBrand)
	api.POST("/brands", CreateBrand)
	api.PUT("/brands/:id", UpdateBrand)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Cfg.HttpPort), router))
}
