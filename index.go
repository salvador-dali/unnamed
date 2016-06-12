package main

import (
	"database/sql"
	"fmt"
	"github.com/dimfeld/httptreemux"
	_ "github.com/lib/pq"
	"github.com/salvador-dali/unnamed/config"
	"github.com/salvador-dali/unnamed/models/brands"
	"log"
	"net/http"
)

// Cfg stores information about all environment variables
var Cfg config.Config

func main() {
	// Create a router, initialize db connection and config
	router, Cfg := httptreemux.New(), config.GetConfig()
	api, dbURL := router.NewGroup("/api/v1"), fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", Cfg.DbUser, Cfg.DbPass, Cfg.DbHost, Cfg.DbPort, Cfg.DbName)
	db, err := sql.Open("postgres", dbURL) // it does not open connection: http://go-database-sql.org/accessing.html
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Brands
	api.GET("/brands", brands.GetAllBrands(db))
	api.GET("/brands/:id", brands.GetBrand(db))
	api.POST("/brands", brands.CreateBrand(db))
	api.PUT("/brands/:id", brands.UpdateBrand)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Cfg.HttpPort), router))
}
