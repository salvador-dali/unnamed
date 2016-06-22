package misc

import (
	"net/mail"
	"strings"
	"time"
)

const (
	passwordMinLen = 8
)

// Error codes
const (
	IdNotNatural                = 100 // ID should be a positive integer
	NameIsNotValid              = 101 // name does not look right. Too long or empty
	WrongNumParams              = 102 // number of parameters is not correct
	FollowYourself              = 103 // user can't follow himself
	TooManyTags                 = 104 // user provided more tags that allowed
	NoTags                      = 105 // user has not provided any tags
	NoSalt                      = 106 // system does not have enough randomness
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

// ErrorCode stores code of a problem that happened while processing client's request.
// It is sent together with 404 status code. It is up to a client how to present it
type ErrorCode struct {
	Id int `json:"error"`
}

// Id stores information about the id of the element which was just inserted
type Id struct {
	Id int `json:"id"`
}

// Brand stores all information about a Brand model
type Brand struct {
	Id        int        `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Issued_at *time.Time `json:"issued_at,omitempty"`
}

// Tag stores all information about a Tag model
type Tag struct {
	Id          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Issued_at   *time.Time `json:"issued_at,omitempty"`
}

// User stores all information about a User model
type User struct {
	Id            int        `json:"id,omitempty"`
	Nickname      string     `json:"nickname,omitempty"`
	Image         string     `json:"image,omitempty"`
	About         string     `json:"about,omitempty"`
	Expertise     int        `json:"expertise,omitempty"`
	Followers_num int        `json:"followers_num,omitempty"`
	Following_num int        `json:"following_num,omitempty"`
	Purchases_num int        `json:"purchases_num,omitempty"`
	Questions_num int        `json:"questions_num,omitempty"`
	Answers_num   int        `json:"answers_num,omitempty"`
	Issued_at     *time.Time `json:"issued_at,omitempty"`
}

// Purchase stores all information about a Purchase model
type Purchase struct {
	Id          int        `json:"id,omitempty"`
	Image       string     `json:"image,omitempty"`
	Description string     `json:"description,omitempty"`
	User_id     int        `json:"user_id,omitempty"`
	Issued_at   *time.Time `json:"issued_at,omitempty"`
	Tags        []int      `json:"tags,omitempty"`
	Brand       int        `json:"brand,omitempty"`
	Likes_num   int        `json:"likes_num,omitempty"`
}

// JwtToken stores authorization information about a user
type JwtToken struct {
	UserId int
	Exp    int
}

func IsPasswordValid(str string) bool {
	return len(str) >= passwordMinLen
}

func IsIdValid(id int) bool {
	return id > 0
}

func ValidateString(str string, maxLen int) (string, bool) {
	str = strings.TrimSpace(str)
	if len(str) == 0 || len(str) > maxLen {
		return "", false
	}

	return str, true
}

func ValidateEmail(str string) (string, bool) {
	e, err := mail.ParseAddress(strings.ToLower(strings.TrimSpace(str)))
	if err != nil {
		return "", false
	}

	return e.Address, true
}
