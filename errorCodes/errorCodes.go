// Package errorCodes contains all possible errors that can be returned to a client
package errorCodes

const (
	IdNotNatural   = 100 // ID should be a positive integer
	IdNotExist     = 101 // No record is found with such ID
	DuplicateName  = 102 // Such name already exists
	NameIsNotValid = 103 // Name does not look right. Too long or empty
)
