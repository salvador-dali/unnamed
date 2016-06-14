// Package errorCodes contains all possible errors that can be returned to a client
package errorCodes

// Errors related to client's input validation
const (
	IdNotNatural      = 100 // ID should be a positive integer
	IdNotExist        = 101 // No record is found with such ID
	DuplicateName     = 102 // Such name already exists
	NameIsNotValid    = 103 // Name does not look right. Too long or empty
	WrongNumParams    = 104 // Number of parameters is not correct
	DuplicateFollower = 105 // You already follow this user
)

// Errors related to a database operations
const (
	DbNothingToReport = 0   // either there was no error, or client should not know about it
	DbValueTooLong    = 600 // text value is too long. Inserted a string of length X + 1 in Varchar(X)
	DbDuplicate       = 601 // duplicate constrain violation. Inserted X, where X already exist and should be unique
	DbNoElement       = 602 // Was searching for an element by ID. Have not found it
)
