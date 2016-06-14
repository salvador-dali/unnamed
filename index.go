package main

import (
	"../unnamed/config"
	"../unnamed/routes"
	"../unnamed/storage"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

func main() {
	// Initializes config and a database
	Cfg := config.Init()
	storage.Init(Cfg.DbUser, Cfg.DbPass, Cfg.DbHost, Cfg.DbName, Cfg.DbPort)
	defer storage.Db.Close()

	// Creates a router
	router := httptreemux.New()
	api := router.NewGroup("/api/v1")

	// Brands
	api.GET("/brands", routes.GetAllBrands)
	api.GET("/brands/:id", routes.GetBrand)
	api.POST("/brands", routes.CreateBrand)
	api.PUT("/brands/:id", routes.UpdateBrand)

	// Tags
	api.GET("/tags", routes.GetAllTags)
	api.GET("/tags/:id", routes.GetTag)
	api.POST("/tags", routes.CreateTag)
	api.PUT("/tags/:id", routes.UpdateTag)

	// Users
	api.GET("/users/:id", routes.GetUser)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Cfg.HttpPort), router))
}
