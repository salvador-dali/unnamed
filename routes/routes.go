// Package routes is responsible for validating clients input and submitting output to a client
package routes

import (
	"../auth"
	"../misc"
	"../models/brand"
	"../models/purchase"
	"../models/tag"
	"../models/user"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// sendJson sends a JSON back to a client with a status Code. Makes error checking
func sendJson(w http.ResponseWriter, data interface{}, statusCode int) {
	if json, err := json.Marshal(data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		w.Write(json)
	}
}

// readJson parses request body and sends BadRequest status if can't be parsed
func readJson(r *http.Request, w http.ResponseWriter) ([]byte, bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return []byte{}, false
	}
	return body, true
}

// isCodeTrivial checks if the code is important enough to be treated as a failure
// NothingToReport is the only code which is unimportant (there might be some failures that were
// written in the log, but a client should not know about them at all
// In all other reasons server sends just a code number to represent a problem. Codes below 200
// represents failure to find something. Was looking by a userId/purchaseID and has found nothing.
func isCodeTrivial(code int, w http.ResponseWriter) bool {
	if code == misc.NothingToReport {
		return true
	}

	if code < 200 {
		sendJson(w, misc.ErrorCode{code}, http.StatusNotFound)
		return false
	}

	sendJson(w, misc.ErrorCode{code}, http.StatusBadRequest)
	return false
}

// validateNaturalNumber checks if the value is a natural and returns it.
// If not, sends a 404 status code and responds with an error JSON
func validateNumeric(w http.ResponseWriter, id string) int {
	id_valid, err := strconv.Atoi(id)
	if err != nil || id_valid <= 0 {
		sendJson(w, misc.ErrorCode{misc.NotNatural}, http.StatusNotFound)
		return 0
	}
	return id_valid
}

// getUserId parses a token header for a JWT token. If found, it is validated and a userId is returned
// otherwise it returns 0. If ResponseWriter is specified, it additionally sends Unauthorized header
func getUserId(r *http.Request, w http.ResponseWriter) int {
	jwt := r.Header.Get("token")
	if len(jwt) == 0 {
		if w != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
		return 0
	}

	jwtToken, err := auth.ValidateJWT(jwt)
	if err != nil {
		if w != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
		return 0
	}

	if !jwtToken.Verified {
		w.WriteHeader(http.StatusUnauthorized)
		return 0
	}

	return jwtToken.UserId
}

// extractPurchasesWithId simplifies extracting many purchases knowing some id
type getPurchasesHelper func(int) ([]*misc.Purchase, int)

func extractPurchasesHelperSendJson(getData getPurchasesHelper, w http.ResponseWriter, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if data, code := getData(id); isCodeTrivial(code, w) {
		sendJson(w, data, http.StatusOK)
	}
}

// GetAllBrands returns all the brands (id, name)
func GetAllBrands(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if brands, code := brand.ShowAll(); isCodeTrivial(code, w) {
		sendJson(w, brands, http.StatusOK)
	}
}

// GetBrand returns full information about a brand
func GetBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if brand, code := brand.ShowById(id); isCodeTrivial(code, w) {
		sendJson(w, brand, http.StatusOK)
	}
}

// CreateBrand creates a brand with a specific name
func CreateBrand(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonName
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if getUserId(r, w) == 0 {
		return
	}

	if id, code := brand.Create(data.Name); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

// UpdateBrand changes the brand's name for a specific brandID
func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	var data misc.JsonName
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if getUserId(r, w) == 0 {
		return
	}

	if code := brand.Update(id, data.Name); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// GetAllTags returns all the tags (id, name)
func GetAllTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if tags, code := tag.ShowAll(); isCodeTrivial(code, w) {
		sendJson(w, tags, http.StatusOK)
	}
}

// GetTag returns full information about a tag
func GetTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if tag, code := tag.ShowById(id); isCodeTrivial(code, w) {
		sendJson(w, tag, http.StatusOK)
	}
}

// CreateTag creates a tag with a specific name and description
func CreateTag(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonNameDescr
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if getUserId(r, w) == 0 {
		return
	}

	if id, code := tag.Create(data.Name, data.Descr); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

// UpdateTag changes the tag's name for a specific tagID
func UpdateTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	var data misc.JsonNameDescr
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if getUserId(r, w) == 0 {
		return
	}

	if code := tag.Update(id, data.Name, data.Descr); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// GetUser returns full information about a user
func GetUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if user, code := user.ShowById(id); isCodeTrivial(code, w) {
		sendJson(w, user, http.StatusOK)
	}
}

// UpdateYourUserInfo changes the information about a user who is currently
func UpdateUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonNicknameAbout
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if code := user.Update(userId, data.Nickname, data.About); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// Follow a current user starts following some user
func Follow(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if code := user.Follow(userId, id); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// Unfollow a current stops following some user
func Unfollow(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if code := user.Unfollow(userId, id); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// GetFollowing returns all the users, whom this user follows
func GetFollowing(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if users, code := user.GetFollowing(id); isCodeTrivial(code, w) {
		sendJson(w, users, http.StatusOK)
	}
}

// GetFollowers returns all the users, who follows this user
func GetFollowers(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if users, code := user.GetFollowers(id); isCodeTrivial(code, w) {
		sendJson(w, users, http.StatusOK)
	}
}

// Login returns a jwt token if a user passed correct credentials
func Login(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonEmailPassword
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if jwt, ok := user.Login(data.Email, data.Password); ok {
		sendJson(w, misc.Jwt{jwt}, http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// ExtendJwt takes a valid JWT token and issues a new one with a full TTL
func ExtendJwt(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if jwt, err := auth.ExtendJWT(r.Header.Get("token")); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		sendJson(w, misc.Jwt{jwt}, http.StatusOK)
	}
}

// CreateUser creates a new unconfirmed user
func CreateUser(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonNicknameEmailPassword
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	if id, code := user.Create(data.Nickname, data.Email, data.Password); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

// VerifyEmail verifies a previously unconfirmed user
func VerifyEmail(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	userId := validateNumeric(w, ps["id"])
	if userId <= 0 {
		return
	}

	if jwt, ok := user.VerifyEmail(userId, ps["code"]); ok {
		sendJson(w, misc.Jwt{jwt}, http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetAllPurchases returns all the purchases in reverse order
func GetAllPurchases(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if purchases, code := purchase.ShowAll(); isCodeTrivial(code, w) {
		sendJson(w, purchases, http.StatusOK)
	}
}

// GetUserPurchases returns all the list of all purchases done by this user in reverse order
func GetUserPurchases(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	extractPurchasesHelperSendJson(purchase.ShowByUserId, w, ps)
}

// GetAllPurchases returns all the purchases which were tagged with a particular brand
func GetAllPurchasesWithBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	extractPurchasesHelperSendJson(purchase.ShowByBrandId, w, ps)
}

// GetAllPurchases returns all the purchases which were tagged with a particular tag
func GetAllPurchasesWithTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	extractPurchasesHelperSendJson(purchase.ShowByTagId, w, ps)
}

// GetPurchase returns full information about a purchase
func GetPurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateNumeric(w, ps["id"])
	if id <= 0 {
		return
	}

	if purchase, code := purchase.ShowById(id); isCodeTrivial(code, w) {
		sendJson(w, purchase, http.StatusOK)
	}
}

// CreatePurchase allows a current user to create a purchase
func CreatePurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	var data misc.JsonDescrBrandTag
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if id, code := purchase.Create(userId, data.Descr, data.BrandId, data.TagIds); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

// LikePurchase allows current user to like a particular purchase
func LikePurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	purchaseId := validateNumeric(w, ps["id"])
	if purchaseId <= 0 {
		return
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if code := purchase.Like(purchaseId, userId); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// UnlikePurchase allows current user to revert his like of a particular purchase
func UnlikePurchase(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	purchaseId := validateNumeric(w, ps["id"])
	if purchaseId <= 0 {
		return
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if code := purchase.Unlike(purchaseId, userId); isCodeTrivial(code, w) {
		sendJson(w, nil, http.StatusNoContent)
	}
}

// AskQuestion allows current user to ask a question about someone's purchase
func AskQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	purchaseId := validateNumeric(w, ps["id"])
	if purchaseId <= 0 {
		return
	}

	var data misc.JsonName
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if id, code := purchase.AskQuestion(purchaseId, userId, data.Name); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

// AnswerQuestion allows current user to answer a question about his own purchase
func AnswerQuestion(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	questionId := validateNumeric(w, ps["id"])
	if questionId <= 0 {
		return
	}

	var data misc.JsonName
	if body, ok := readJson(r, w); !ok {
		return
	} else {
		json.Unmarshal(body, &data)
	}

	userId := getUserId(r, w)
	if userId == 0 {
		return
	}

	if id, code := purchase.AnswerQuestion(questionId, userId, data.Name); isCodeTrivial(code, w) {
		sendJson(w, misc.Id{id}, http.StatusCreated)
	}
}

func Avatar(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	misc.SaveFileFromClient(w, r, "img")
}
