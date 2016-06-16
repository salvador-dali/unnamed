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
	api.PUT("/users/me/info", routes.UpdateYourUserInfo)
	api.POST("/users/me/follow/:id", routes.Follow)
	api.DELETE("/users/me/follow/:id", routes.Unfollow)
	api.GET("/users/:id/followers", routes.GetFollowers)
	api.GET("/users/:id/following", routes.GetFollowing)
	api.GET("/users/:id/purchases", routes.GetUserPurchases)
	api.GET("/purchases", routes.GetAllPurchases)
	api.GET("/purchases/brand/:id", routes.GetAllPurchasesWithBrand)
	api.GET("/purchases/tag/:id", routes.GetAllPurchasesWithTag)
	api.GET("/purchases/:id", routes.GetPurchase)
	api.POST("/purchases/:id/like", routes.LikePurchase)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Cfg.HttpPort), router))
}
