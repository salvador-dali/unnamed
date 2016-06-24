package main

import (
	"./config"
	"./routes"
	"./storage"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

func main() {
	// Initializes config and a database
	config.Init()
	storage.Init(config.Cfg.DbUser, config.Cfg.DbPass, config.Cfg.DbHost, config.Cfg.DbName, config.Cfg.DbPort)
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
	//api.GET("/users/:id/purchases", routes.GetUserPurchases)
	//
	//// Purchases
	//api.GET("/purchases", routes.GetAllPurchases)
	//api.POST("/purchases", routes.CreatePurchase)
	//api.GET("/purchases/brand/:id", routes.GetAllPurchasesWithBrand)
	//api.GET("/purchases/tag/:id", routes.GetAllPurchasesWithTag)
	//api.GET("/purchases/:id", routes.GetPurchase)
	//api.POST("/purchases/:id/like", routes.LikePurchase)
	//api.DELETE("/purchases/:id/like", routes.UnlikePurchase)
	//api.POST("/purchases/:id/ask", routes.AskQuestion)
	//
	//// Questions
	////api.POST("/questions/:id/vote", routes.UpvoteQuestion)
	////api.DELETE("/questions/:id/vote", routes.DownvoteQuestion)
	//api.POST("/questions/:id/answer", routes.AnswerQuestion)

	// Answers
	//api.POST("/answer/:id/vote", routes.UpvoteAnswer)
	//api.DELETE("/answer/:id/vote", routes.DownvoteAnswer)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.HttpPort), router))
}
