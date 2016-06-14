// Package errorCodes contains all possible errors that can be returned to a client
package errorCodes

// Errors related to client's input validation
const (
	IdNotNatural      = 100 // ID should be a positive integer
	NameIsNotValid    = 101 // name does not look right. Too long or empty
	WrongNumParams    = 102 // number of parameters is not correct
)

// Errors related to a database operations
const (
	DbNothingToReport = 0   // either there was no error, or client should not know about it
	DbValueTooLong    = 600 // text value is too long. Inserted a string of length X + 1 in Varchar(X)
	DbDuplicate       = 601 // duplicate constrain violation. Inserted X, where X already exist and should be unique
	DbNoElement       = 602 // was searching for an element by ID. Have not found it
)
