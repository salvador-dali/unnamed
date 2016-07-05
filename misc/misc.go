package misc

import (
	"math/rand"
	"net/mail"
	"strings"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	passwordMinLen = 8
	ConfCodeLen    = 20   // length of the confirmation code which will be sent to a newly created user
	MaxTags        = 4    // maximum number of tags possible for a purchase
	MaxLenS        = 40   // maximum length of the small field in SQL
	MaxLenB        = 1000 // maximum length of the big field in SQL
)

// Error codes
const (
	NothingToReport = 0   // either there is no error, or a client should not know about it
	NothingUpdated  = 100 // wanted to update an element by ID. Element does not exist
	NoElement       = 101 // searched for an element by ID. Have not found it.
	NoPurchase      = 102 // purchase with such ID does not exist
	NotNatural      = 103 // provided value was not a natural number

	WrongName           = 201 // name is too long or empty
	WrongDescr          = 202 // description is too long or empty
	WrongEmail          = 203 // email does not look right
	WrongPassword       = 204 // password is too short
	WrongTagsNum        = 205 // user provided more tags that allowed
	WrongTags           = 206 // a one or more tags are not in the database
	FollowYourself      = 207 // user can't follow himself
	VoteForYourself     = 208 // a person should not vote for his own stuff
	AskYourself         = 209 // a person should not ask questions about his purchase
	AnswerOtherPurchase = 210 // user can answer only question about his purchase
	NoTags              = 211 // user has not provided any tags
	WrongImg            = 212 // something wrong with the image

	NoSalt                = 301 // system does not have enough randomness
	DbDuplicate           = 302 // duplicate constrain violation. Inserted X, where X already exists and should be unique
	DbForeignKeyViolation = 303 // foreign key violation
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

// Id stores jwt token
type Jwt struct {
	Jwt string `json:"token"`
}

// Brand stores all information about a Brand model
type Brand struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Issued_at int64  `json:"issued_at,omitempty"`
}

// Tag stores all information about a Tag model
type Tag struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Issued_at   int64  `json:"issued_at,omitempty"`
}

// User stores all information about a User model
type User struct {
	Id            int    `json:"id,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	Image         string `json:"image,omitempty"`
	About         string `json:"about,omitempty"`
	Expertise     int    `json:"expertise,omitempty"`
	Followers_num int    `json:"followers_num,omitempty"`
	Following_num int    `json:"following_num,omitempty"`
	Purchases_num int    `json:"purchases_num,omitempty"`
	Questions_num int    `json:"questions_num,omitempty"`
	Answers_num   int    `json:"answers_num,omitempty"`
	Issued_at     int64  `json:"issued_at,omitempty"`
}

// Purchase stores all information about a Purchase model
type Purchase struct {
	Id          int    `json:"id,omitempty"`
	Image       string `json:"image,omitempty"`
	Description string `json:"description,omitempty"`
	User_id     int    `json:"user_id,omitempty"`
	Issued_at   int64  `json:"issued_at,omitempty"`
	Tags        []int  `json:"tags,omitempty"`
	Brand       int    `json:"brand,omitempty"`
	Likes_num   int    `json:"likes_num,omitempty"`
}

// JwtToken stores authorization information about a user
type JwtToken struct {
	UserId   int
	Iat      int
	Exp      int
	Verified bool
}

type JsonName struct {
	Name string `json:"name"`
}

type JsonNameDescr struct {
	Name  string `json:"name"`
	Descr string `json:"descr"`
}

type JsonNicknameAbout struct {
	Nickname string `json:"nickname"`
	About    string `json:"about"`
}

type JsonDescrBrandTag struct {
	Descr   string `json:"descr"`
	BrandId int    `json:"brand"`
	TagIds  []int  `json:"tags"`
}

type JsonEmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JsonNicknameEmailPassword struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
