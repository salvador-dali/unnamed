package main

import (
	"./config"
	"./mailer"
	"./psql"
	"./routes"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Init prepares the service for a work:
// - initializes randomness
// - creates a config
// - creates a mailer object
// - creates a database connection
func Init() {
	rand.Seed(time.Now().UnixNano())
	config.Init()
	mailer.Init()
	psql.Init()
}

func main() {
	Init()

	// Creates a router
	router := httptreemux.New()
	api := router.NewGroup("/api/v1")

	// Image
	api.POST("/image/avatar", routes.UploadImageAvatar)
	api.POST("/image/purchase", routes.UploadImagePurchase)

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
	api.POST("/users/login", routes.Login)
	api.GET("/users/login/extend", routes.ExtendJwt)
	api.POST("/users", routes.CreateUser)
	api.GET("/users/:id", routes.GetUser)
	api.PUT("/users/me/info", routes.UpdateUser)
	api.POST("/users/me/follow/:id", routes.Follow)
	api.DELETE("/users/me/follow/:id", routes.Unfollow)
	api.GET("/users/:id/followers", routes.GetFollowers)
	api.GET("/users/:id/following", routes.GetFollowing)
	api.GET("/users/:id/purchases", routes.GetUserPurchases)
	api.GET("/users/verify/:id/:code", routes.VerifyEmail)

	// Purchases
	api.GET("/purchases", routes.GetAllPurchases)
	api.POST("/purchases", routes.CreatePurchase)
	api.GET("/purchases/brand/:id", routes.GetAllPurchasesWithBrand)
	api.GET("/purchases/tag/:id", routes.GetAllPurchasesWithTag)
	api.GET("/purchases/:id", routes.GetPurchase)
	api.POST("/purchases/:id/like", routes.LikePurchase)
	api.DELETE("/purchases/:id/like", routes.UnlikePurchase)
	api.POST("/purchases/:id/ask", routes.AskQuestion)

	// Questions
	//api.POST("/questions/:id/vote", routes.UpvoteQuestion)
	//api.DELETE("/questions/:id/vote", routes.DownvoteQuestion)
	api.POST("/questions/:id/answer", routes.AnswerQuestion)

	// Answers
	//api.POST("/answer/:id/vote", routes.UpvoteAnswer)
	//api.DELETE("/answer/:id/vote", routes.DownvoteAnswer)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.HttpPort), router))
}
