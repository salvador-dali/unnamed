// Package routes is responsible for validating clients input and submitting output to a client
package routes

import (
	"../../unnamed/errorCodes"
	"../../unnamed/storage"
	"../../unnamed/structs"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// -- A couple of Helper methods and constants
const (
	maximumNameLength  = 40
	maximumDescrLength = 1000
)

// sendJSON sends a JSON back to a client with a status Code. Makes error checking
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(statusCode)
	w.Write(json)
}

// isErrorReasonSerious checks whether an error happened and whether the reason is serious to notify
// a client. If there is an error and no reason to notify a client - just halt the operation and write to log
// if there is a reason - just write a JSON with an error code to a client
func isErrorReasonSerious(err error, reason int, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	if reason == errorCodes.DbNoElement {
		sendJSON(w, structs.ErrorCode{reason}, http.StatusNotFound)
		return true
	} else if reason > 0 {
		sendJSON(w, structs.ErrorCode{reason}, http.StatusBadRequest)
		return true
	}

	log.Fatal(err)
	return true
}

// validateId checks whether ID is a natural number and returns it.
// If not, sends a 404 status code and responds with an error JSON
func validateId(w http.ResponseWriter, id string) int {
	id_valid, err := strconv.Atoi(id)
	if err != nil || id_valid <= 0 {
		sendJSON(w, structs.ErrorCode{errorCodes.IdNotNatural}, http.StatusNotFound)
		return 0
	}
	return id_valid
}

// validateName checks whether Name is not empty and has a correct length.
// If not, sends a 404 status code and responds with an error JSON
func validateName(w http.ResponseWriter, name string, maxLen int) (string, bool) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > maxLen {
		sendJSON(w, structs.ErrorCode{errorCodes.NameIsNotValid}, http.StatusBadRequest)
		return "", false
	}

	return name, true
}

// isValidFormLength returns true if exactly validLen parameters were passed
// If not, sends a 404 status code and responds with an error JSON
func isValidFormLength(w http.ResponseWriter, r *http.Request, validLen int) bool {
	r.ParseForm()
	if len(r.Form) == validLen {
		return true
	}

	sendJSON(w, structs.ErrorCode{errorCodes.WrongNumParams}, http.StatusBadRequest)
	return false
}

// -- Actual Handlers

// GetAllBrands returns all the brands (id, name)
func GetAllBrands(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	brands, err, reason := storage.GetAllBrands()
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, brands, http.StatusOK)
}

// GetBrand returns full information about a brand
func GetBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateId(w, ps["id"])
	if id <= 0 {
		return
	}

	brand, err, reason := storage.GetBrand(id)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, brand, http.StatusOK)
}

// CreateBrand creates a brand with a specific name
func CreateBrand(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")
	if !isValidFormLength(w, r, 1) {
		return
	}

	name, ok := validateName(w, r.PostFormValue("name"), maximumNameLength)
	if !ok {
		return
	}

	id, err, reason := storage.CreateBrand(name)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, structs.Id{int(id)}, http.StatusCreated)
}

// UpdateBrand changes the brand's name for a specific brandID
func UpdateBrand(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if !isValidFormLength(w, r, 1) {
		return
	}

	id := validateId(w, ps["id"])
	if id <= 0 {
		return
	}

	name, ok := validateName(w, r.PostFormValue("name"), maximumNameLength)
	if !ok {
		return
	}

	err, reason := storage.UpdateBrand(id, name)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, nil, http.StatusNoContent)
}

// GetAllTags returns all the tags (id, name)
func GetAllTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	tags, err, reason := storage.GetAllTags()
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, tags, http.StatusOK)
}

// GetTag returns full information about a tag
func GetTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateId(w, ps["id"])
	if id <= 0 {
		return
	}

	tag, err, reason := storage.GetTag(id)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, tag, http.StatusOK)
}

// CreateTag creates a tag with a specific name and description
func CreateTag(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")
	if !isValidFormLength(w, r, 2) {
		return
	}

	name, ok := validateName(w, r.PostFormValue("name"), maximumNameLength)
	if !ok {
		return
	}

	descr, ok := validateName(w, r.PostFormValue("description"), maximumDescrLength)
	if !ok {
		return
	}

	id, err, reason := storage.CreateTag(name, descr)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, structs.Id{int(id)}, http.StatusCreated)
}

// UpdateTag changes the tag's name for a specific tagID
func UpdateTag(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	if !isValidFormLength(w, r, 2) {
		return
	}

	id := validateId(w, ps["id"])
	if id <= 0 {
		return
	}

	name, ok := validateName(w, r.PostFormValue("name"), maximumNameLength)
	if !ok {
		return
	}

	descr, ok := validateName(w, r.PostFormValue("description"), maximumDescrLength)
	if !ok {
		return
	}

	err, reason := storage.UpdateTag(id, name, descr)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, nil, http.StatusNoContent)
}

// GetUser returns full information about a user
func GetUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Header().Set("Content-Type", "application/javascript")

	id := validateId(w, ps["id"])
	if id <= 0 {
		return
	}

	user, err, reason := storage.GetUser(id)
	if isErrorReasonSerious(err, reason, w) {
		return
	}

	sendJSON(w, user, http.StatusOK)
}
