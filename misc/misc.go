package misc

import (
	"net/mail"
	"strings"
	"time"
)

const (
	passwordMinLen = 8
	MaxTags        = 4		// maximum number of tags possible for a purchase
	MaxLenS        = 40		// maximum length of the small field in SQL
	MaxLenB        = 1000	// maximum length of the big field in SQL
)

// Error codes
const (
	NothingToReport       = 0   // either there is no error, or a client should not know about it
	NothingUpdated        = 100 // wanted to update an element by ID. Element does not exist
	WrongName             = 101 // name is too long or empty
	WrongDescr            = 102 // description is too long or empty
	WrongEmail            = 802 // email does not look right
	WrongPassword         = 803 // password is too short
	WrongParamsNum        = 103 // number of parameters is not correct
	WrongTagsNum          = 104 // user provided more tags that allowed
	WrongTags             = 105 // a one or more tags are not in the database
	FollowYourself        = 106 // user can't follow himself
	VoteForYourself       = 107 // a person should not vote for his own stuff
	AskYourself           = 108 // a person should not ask questions about his purchase
	AnswerOtherPurchase   = 109 // user can answer only question about his purchase
	NoTags                = 110 // user has not provided any tags
	NoSalt                = 111 // system does not have enough randomness
	NoElement             = 112 // searched for an element by ID. Have not found it.
	NoPurchase            = 113 // purchase with such ID does not exist
	DbDuplicate           = 114 // duplicate constrain violation. Inserted X, where X already exists and should be unique
	DbForeignKeyViolation = 115 // foreign key violation
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
