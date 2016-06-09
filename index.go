package main

import (
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

 //----------- Brands
func GetAllBrands(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func CreateBrand(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ----------- Tags
func GetAllTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func CreateTag(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ----------- Users
func CreateUser(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func LoginUser(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func LogoutUser(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func VerifyUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UpdateUserAvatar(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func FollowUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetFollowers(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetFollowing(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetPurchases(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ----------- Purchases
func GetAllPurchases(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetAllPurchasesWithTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetAllPurchasesWithBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func CreatePurchase(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetPurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func LikePurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UnlikePurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func AskQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ----------- Questions
func UpvoteQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func DownvoteQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func AnswerQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// ----------- Answers
func UpvoteAnswer(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func DownvoteAnswer(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	router := httptreemux.New()
	api := router.NewGroup("/api/v1")

	// Brands
	api.GET("/brands/", GetAllBrands)
	api.GET("/brands/:id/", GetBrand)
	api.POST("/brands/", CreateBrand)
	api.PUT("/brands/:id/", UpdateBrand)

	// Tags
	api.GET("/tags/", GetAllTags)
	api.GET("/tags/:id/", GetTag)
	api.POST("/tags/", CreateTag)
	api.PUT("/tags/:id/", UpdateTag)

	// Users
	api.POST("/users/", CreateUser)
	api.POST("/users/login/", LoginUser)
	api.POST("/users/logout/", LogoutUser)
	api.GET("/users/:id/", GetUser)
	api.GET("/users/me/email/:hash/", VerifyUser)
	api.PUT("/users/me/info/", UpdateUserInfo)
	api.PUT("/users/me/avatar/", UpdateUserAvatar)
	api.POST("/users/me/follow/:id/", FollowUser)
	api.DELETE("/users/me/follow/:id/", UnfollowUser)
	api.GET("/users/:id/followers/", GetFollowers)
	api.GET("/users/:id/following/", GetFollowing)
	api.GET("/users/:id/purchases/", GetPurchases)

	// Purchases
	api.GET("/purchases/", GetAllPurchases)
	api.GET("/purchases/tag/:id/", GetAllPurchasesWithTag)
	api.GET("/purchases/brand/:id/", GetAllPurchasesWithBrand)
	api.POST("/purchases/", CreatePurchase)
	api.GET("/purchases/:id/", GetPurchase)
	api.POST("/purchases/:id/like/", LikePurchase)
	api.DELETE("/purchases/:id/like/", UnlikePurchase)
	api.POST("/purchases/:id/ask/", AskQuestion)

	// Questions
	api.POST("/questions/:id/vote/", UpvoteQuestion)
	api.DELETE("/questions/:id/vote/", DownvoteQuestion)
	api.POST("/questions/:id/answer/", AnswerQuestion)

	// Answers
	api.POST("/answer/:id/vote/", UpvoteAnswer)
	api.DELETE("/answer/:id/vote/", DownvoteAnswer)

    log.Fatal(http.ListenAndServe(":8080", router))
}
