// Package errorCodes contains all possible errors that can be returned to a client
package errorCodes

// Errors related to client's input validation
const (
	IdNotNatural   = 100 // ID should be a positive integer
	NameIsNotValid = 101 // name does not look right. Too long or empty
	WrongNumParams = 102 // number of parameters is not correct
	FollowYourself = 103 // user can't follow himself
	TooManyTags    = 104 // user provided more tags that allowed
	NoTags         = 105 // user has not provided any tags
)

// Errors related to a database operations
const (
	DbNothingToReport           = 0   // either there is no error, or a client should not know about it
	DbValueTooLong              = 600 // text value is too long. Tried to insert a string > X in Varchar(X)
	DbDuplicate                 = 601 // duplicate constrain violation. Inserted X, where X already exists and should be unique
	DbNoElement                 = 602 // searched for an element by ID. Have not found it
	DbNothingUpdated            = 603 // wanted to update an element by ID. Element does not exist
	DbNoPurchase                = 604 // purchase with such ID does not exist
	DbVoteForOwnStuff           = 605 // a person should not vote for his own stuff
	DbAskAboutOwnStuff          = 606 // a person should not ask questions about his purchase
	DbNotAllTagsCorrect         = 607 // a one or more tags are not in the database
	DbNoPurchaseForQuestion     = 608 // purchase for this question ID does not exist
	DbCannotAnswerOtherPurchase = 609 // user can answer only question about his purchase
	DbForeignKeyViolation       = 610 // foreign key violation
)
