// Package helpers contains various methods that validate fields and send responses to
// a client when errors are found
package structs

import (
	"encoding/json"
	"github.com/salvador-dali/unnamed/errorCodes"
	"github.com/salvador-dali/unnamed/structs"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// validateId checks whether ID is a natural number and returns it.
// If not, sends a 404 status code and responds with an error JSON
func ValidateId(w http.ResponseWriter, id string) (int, bool) {
	id_valid, err := strconv.Atoi(id)
	if err != nil || id_valid <= 0 {
		json, err := json.Marshal(structs.ErrorCode{errorCodes.IdNotNatural})
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		return 0, false
	}
	return id_valid, true
}

// ValidateName checks whether Name is not empty and has a correct length.
// If not, sends a 404 status code and responds with an error JSON
func ValidateName(w http.ResponseWriter, name string, maxLen int) (string, bool) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > maxLen {
		json, err := json.Marshal(structs.ErrorCode{errorCodes.NameIsNotValid})
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
		return "", false
	}

	return name, true
}

// IsValidateFormLength returns true if exactly validLen parameters were passed
// If not, sends a 404 status code and responds with an error JSON
func IsValidateFormLength(w http.ResponseWriter, r *http.Request, validLen int) bool {
	r.ParseForm()
	if len(r.Form) == validLen {
		return true
	}

	json, err := json.Marshal(structs.ErrorCode{errorCodes.WrongNumParams})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(json)
	return false
}
