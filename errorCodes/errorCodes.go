// Package errorCodes contains all possible errors that can be returned to a client
package errorCodes

// Errors related to client's input validation
const (
	IdNotNatural   = 100 // ID should be a positive integer
	NameIsNotValid = 101 // name does not look right. Too long or empty
	WrongNumParams = 102 // number of parameters is not correct
	FollowYourself = 103 // you are trying to follow yourself
	TooManyTags    = 104 // maximum number of tags for a purchase is X, you provides > X
	NoTags         = 105 // no tags were provided
)

// Errors related to a database operations
const (
	DbNothingToReport   = 0   // either there was no error, or client should not know about it
	DbValueTooLong      = 600 // text value is too long. Inserted a string of length X + 1 in Varchar(X)
	DbDuplicate         = 601 // duplicate constrain violation. Inserted X, where X already exist and should be unique
	DbNoElement         = 602 // was searching for an element by ID. Have not found it
	DbNothingUpdated    = 603 // was trying to update an element, but the element was not found
	DbNoPurchase        = 604 // purchase with such ID does not exist
	DbVoteForOwnStuff   = 605 // a person should not vote for his own stuff
	DbAskAboutOwnStuff  = 606 // a person should not ask questions about his purchase
	DbNotAllTagsCorrect = 607 // a one or more tags are not in the database
)
