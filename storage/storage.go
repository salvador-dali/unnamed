// Package storage is responsible for all database operations
package storage

import (
	"../auth"
	"../misc"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var Db *sql.DB

// Init prepares the database abstraction for later use. It does not establish any connections to
// the database, nor does it validate driver connection parameters. To do this call Ping
// http://go-database-sql.org/accessing.html
func Init(user, pass, host, name string, port int) {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, name)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	Db = db
}

// a couple of helper database functions

// isAffectedOneRow checks that the result of a pure INSERT or UPDATE executed with Exec has modified only 1 row
func isAffectedOneRow(sqlResult sql.Result) (error, int) {
	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		return err, misc.NothingToReport
	}

	if affectedRows == 1 {
		return nil, misc.NothingToReport
	} else if affectedRows == 0 {
		return errors.New("nothing updated"), misc.NothingUpdated
	}
	return errors.New(fmt.Sprintf("Expected to update 1 value. %d updated", affectedRows)), misc.NothingToReport
}

// checkSpecificDriverErrors analyses the error result against specific errors that a client should know about
// This checks for Value Limit violation and Duplicate constraint violation
func checkSpecificDriverErrors(err error) (error, int) {
	if errPg, ok := err.(*pq.Error); ok {
		s := string(errPg.Code)
		switch {
		case s == "23505":
			return err, misc.DbDuplicate
		case s == "22001":
			return err, misc.WrongName
		case s == "23503":
			return err, misc.DbForeignKeyViolation
		}
	}

	return err, misc.NothingToReport
}

// --- Brands ---

func GetAllBrands() ([]*misc.Brand, error, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM brands
		WHERE id > 0`)
	if err != nil {
		return []*misc.Brand{}, err, misc.NothingToReport
	}
	defer rows.Close()

	brands := []*misc.Brand{}
	for rows.Next() {
		brand := misc.Brand{}
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			return []*misc.Brand{}, err, misc.NothingToReport
		}
		brands = append(brands, &brand)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Brand{}, err, misc.NothingToReport
	}

	return brands, nil, misc.NothingToReport
}

func GetBrand(brandId int) (misc.Brand, error, int) {
	if !misc.IsIdValid(brandId) {
		return misc.Brand{}, errors.New("Not positive id"), misc.NoElement
	}

	brand := misc.Brand{}
	if err := Db.QueryRow(`
		SELECT name, issued_at
		FROM brands
		WHERE id = $1`, brandId,
	).Scan(&brand.Name, &brand.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return misc.Brand{}, err, misc.NoElement
		}

		return misc.Brand{}, err, misc.NothingToReport
	}

	brand.Id = brandId
	return brand, nil, misc.NothingToReport
}

func CreateBrand(name string) (int, error, int) {
	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		return 0, errors.New("Wrong name"), misc.WrongName
	}

	brandId := 0
	err := Db.QueryRow(`
		INSERT INTO brands (name)
		VALUES ($1)
		RETURNING id`, name,
	).Scan(&brandId)
	if err == nil {
		return brandId, nil, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateBrand(brandId int, name string) (error, int) {
	if !misc.IsIdValid(brandId) {
		return errors.New("Nothing updated"), misc.NothingUpdated
	}

	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		return errors.New("Wrong name"), misc.WrongName
	}

	sqlResult, err := Db.Exec(`
		UPDATE brands
		SET name = $1
		WHERE id = $2`, name, brandId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

// --- Tags ---

func GetAllTags() ([]*misc.Tag, error, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM tags`)
	if err != nil {
		return []*misc.Tag{}, err, misc.NothingToReport
	}
	defer rows.Close()

	tags := []*misc.Tag{}
	for rows.Next() {
		tag := misc.Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return []*misc.Tag{}, err, misc.NothingToReport
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Tag{}, err, misc.NothingToReport
	}

	return tags, nil, misc.NothingToReport
}

func GetTag(tagId int) (misc.Tag, error, int) {
	if !misc.IsIdValid(tagId) {
		return misc.Tag{}, errors.New("Not positive id"), misc.NoElement
	}

	tag := misc.Tag{}
	if err := Db.QueryRow(`
		SELECT name, description, issued_at
		FROM tags
		WHERE id = $1`, tagId,
	).Scan(&tag.Name, &tag.Description, &tag.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return misc.Tag{}, err, misc.NoElement
		}

		return misc.Tag{}, err, misc.NothingToReport
	}

	tag.Id = tagId
	return tag, nil, misc.NothingToReport
}

func CreateTag(name, descr string) (int, error, int) {
	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		return 0, errors.New("Wrong name"), misc.WrongName
	}

	descr, ok = misc.ValidateString(descr, misc.MaxLenB)
	if !ok {
		return 0, errors.New("Wrong description"), misc.WrongDescr
	}

	tagId := 0
	err := Db.QueryRow(`
		INSERT INTO tags (name, description)
		VALUES ($1, $2)
		RETURNING id`, name, descr,
	).Scan(&tagId)
	if err == nil {
		return tagId, nil, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateTag(tagId int, name, descr string) (error, int) {
	if !misc.IsIdValid(tagId) {
		return errors.New("Nothing updated"), misc.NothingUpdated
	}

	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		return errors.New("Wrong name"), misc.WrongName
	}

	descr, ok = misc.ValidateString(descr, misc.MaxLenB)
	if !ok {
		return errors.New("Wrong description"), misc.WrongDescr
	}

	sqlResult, err := Db.Exec(`
		UPDATE tags
		SET name = $1, description = $2
		WHERE id = $3`, name, descr, tagId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

// validateTags makes sure that all the tags are in the database. psql does not support this http://dba.stackexchange.com/q/60132/15318
func validateTags(tagIds []int) (error, int) {
	if len(tagIds) == 0 {
		return errors.New("no tags"), misc.NoTags
	}

	if len(tagIds) > misc.MaxTags {
		return errors.New("too many tags"), misc.WrongTagsNum
	}

	for _, v := range tagIds {
		if v <= 0 {
			return errors.New("tag is negative"), misc.WrongTags
		}
	}

	stringTagIds, num := make([]string, len(tagIds), len(tagIds)), 0
	for k, v := range tagIds {
		stringTagIds[k] = strconv.Itoa(v)
	}

	if err := Db.QueryRow(`
		SELECT COUNT(id)
		FROM tags
		WHERE id IN ($1)`, strings.Join(stringTagIds, ","),
	).Scan(&num); err != nil {
		return err, misc.NothingToReport
	}

	if num != len(tagIds) {
		return errors.New("some tags are missing"), misc.WrongTags
	}

	return nil, misc.NothingToReport
}

// --- Users ---

func GetUser(userId int) (misc.User, error, int) {
	if !misc.IsIdValid(userId) {
		return misc.User{}, errors.New("Not positive id"), misc.NoElement
	}

	user := misc.User{}
	if err := Db.QueryRow(`
		SELECT nickname, image, about, expertise, followers_num, following_num, purchases_num, questions_num, answers_num, issued_at
		FROM users WHERE id = $1`, userId,
	).Scan(
		&user.Nickname, &user.Image, &user.About, &user.Expertise, &user.Followers_num,
		&user.Following_num, &user.Purchases_num, &user.Questions_num, &user.Answers_num,
		&user.Issued_at,
	); err != nil {
		if err == sql.ErrNoRows {
			return misc.User{}, err, misc.NoElement
		}

		return misc.User{}, err, misc.NothingToReport
	}
	user.Id = userId
	return user, nil, misc.NothingToReport
}

func UpdateUser(userId int, nickname, about string) (error, int) {
	if !misc.IsIdValid(userId) {
		return errors.New("Nothing updated"), misc.NothingUpdated
	}

	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		return errors.New("Wrong name"), misc.WrongName
	}

	about, ok = misc.ValidateString(about, misc.MaxLenB)
	if !ok {
		return errors.New("Wrong about"), misc.WrongDescr
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET nickname = $1, about = $2
		WHERE id = $3`, nickname, about, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

func Follow(whoId, whomId int) (error, int) {
	if !misc.IsIdValid(whomId) {
		return errors.New("Not positive id"), misc.NoElement
	}

	if whoId == whomId {
		return errors.New("can't follow yourself"), misc.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		INSERT INTO followers (who_id, whom_id)
		VALUES ($1, $2)`, whoId, whomId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num + 1
		WHERE id = $1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num + 1
		WHERE id = $1`, whoId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	return nil, misc.NothingToReport
}

func Unfollow(whoId, whomId int) (error, int) {
	if !misc.IsIdValid(whomId) {
		return errors.New("Not positive id"), misc.NoElement
	}

	if whoId == whomId {
		return errors.New("can't follow yourself"), misc.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM followers
		WHERE who_id = $1 AND whom_id = $2`, whoId, whomId)
	if err != nil {
		return err, misc.NothingToReport
	}

	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num - 1
		WHERE id = $1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num - 1
		WHERE id = $1`, whoId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	return nil, misc.NothingToReport
}

func GetFollowing(userId int) ([]*misc.User, error, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, nil, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT whom_id
			FROM followers
			WHERE who_id = $1
		)`, userId)
	if err != nil {
		return []*misc.User{}, err, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*misc.User{}, err, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*misc.User{}, err, misc.NothingToReport
	}

	return users, nil, misc.NothingToReport
}

func GetFollowers(userId int) ([]*misc.User, error, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, nil, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT who_id
			FROM followers
			WHERE whom_id = $1
		)`, userId)
	if err != nil {
		return []*misc.User{}, err, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*misc.User{}, err, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*misc.User{}, err, misc.NothingToReport
	}

	return users, nil, misc.NothingToReport
}

func CreateUser(nickname, email, password string) (int, error, int) {
	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		return 0, errors.New("Wrong name"), misc.WrongName
	}

	email, ok = misc.ValidateEmail(email)
	if !ok {
		return 0, errors.New("Wrong email"), misc.WrongEmail
	}

	if !misc.IsPasswordValid(password) {
		return 0, errors.New("Wrong password"), misc.WrongPassword
	}

	salt, err := auth.GenerateSalt()
	if err != nil {
		return 0, err, misc.NoSalt
	}

	hash, err := auth.PasswordHash(password, salt)
	if err != nil {
		return 0, err, misc.NothingToReport
	}

	userId := 0
	err = Db.QueryRow(`
		INSERT INTO users (nickname, email, password, salt)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, nickname, email, hash, salt,
	).Scan(&userId)
	if err == nil {
		return userId, nil, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func Login(email, password string) (string, bool) {
	email, ok := misc.ValidateEmail(email)
	if !ok || !misc.IsPasswordValid(password){
		return "", false
	}

	userId, hash, salt := 0, make([]byte, 32), make([]byte, 16)

	if err := Db.QueryRow(`
		SELECT id, password, salt
		FROM users
		WHERE email = $1`, email,
	).Scan(&userId, &hash, &salt); err != nil {
		return "", false
	}

	hashAttempt, err := auth.PasswordHash(password, salt)
	if err != nil {
		return "", false
	}

	if !reflect.DeepEqual(hashAttempt, hash) {
		return "", false
	}

	jwt, err := auth.CreateJWT(userId)
	if err != nil {
		return "", false
	}
	return jwt, true
}

// --- Purchases ---

func getPurchases(rows *sql.Rows, err error) ([]*misc.Purchase, error, int) {
	if err != nil {
		return []*misc.Purchase{}, err, misc.NothingToReport
	}
	defer rows.Close()

	purchases := []*misc.Purchase{}
	for rows.Next() {
		p := misc.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			return []*misc.Purchase{}, err, misc.NothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Purchase{}, err, misc.NothingToReport
	}

	return purchases, nil, misc.NothingToReport
}

func whoCreatedPurchaseByPurchaseId(purchaseId int) (int, error, int) {
	if !misc.IsIdValid(purchaseId) {
		return 0, errors.New("No purchase"), misc.NoPurchase
	}

	whosePurchase := 0
	if err := Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			return 0, err, misc.NoPurchase
		}

		return 0, err, misc.NothingToReport
	}

	return whosePurchase, nil, misc.NothingToReport
}

func whoCreatedPurchaseByQuestionId(questionId int) (int, error, int) {
	if !misc.IsIdValid(questionId) {
		// if question does not exist, surely there is no purchase for this question
		return 0, errors.New("No purchase"), misc.NoPurchase
	}

	whosePurchase := 0
	if err := Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = (
			SELECT purchase_id
			FROM questions
			WHERE id = $1
		)`, questionId).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			return 0, err, misc.NoPurchase
		}

		return 0, err, misc.NothingToReport
	}

	return whosePurchase, nil, misc.NothingToReport
}

func GetUserPurchases(userId int) ([]*misc.Purchase, error, int) {
	// userId is the current user and is always valid
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE user_id = $1
		ORDER BY issued_at DESC`, userId)

	return getPurchases(rows, err)
}

func CreatePurchase(userId int, description string, brandId int, tagsId []int) (int, error, int) {
	// userID is the current user and should be valid
	description, ok := misc.ValidateString(description, misc.MaxLenB)
	if !ok {
		return 0, errors.New("Wrong description"), misc.WrongDescr
	}

	if brandId < 0 {
		return 0, errors.New("Not positive id"), misc.NoElement
	}

	if err, code := validateTags(tagsId); err != nil {
		return 0, err, code
	}

	stringTagIds, id := make([]string, len(tagsId), len(tagsId)), 0
	for k, v := range tagsId {
		stringTagIds[k] = strconv.Itoa(v)
	}

	tagsToInsert := "{" + strings.Join(stringTagIds, ",") + "}"
	err := Db.QueryRow(`
		INSERT INTO purchases (image, description, user_id, tag_ids, brand_id)
		VALUES ('', $1, $2, $3, $4)
		RETURNING id`, description, userId, tagsToInsert, brandId).Scan(&id)
	if err != nil {
		return 0, err, misc.NothingToReport
	}

	if err, code := checkSpecificDriverErrors(err); err != nil {
		return 0, err, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET purchases_num = purchases_num + 1
		WHERE id=$1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, err, code
	}

	return id, nil, misc.NothingToReport
}

func GetAllPurchases() ([]*misc.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithBrand(brandId int) ([]*misc.Purchase, error, int) {
	if !misc.IsIdValid(brandId) {
		return []*misc.Purchase{}, nil, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE brand_id = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithTag(tagId int) ([]*misc.Purchase, error, int) {
	if !misc.IsIdValid(tagId) {
		return []*misc.Purchase{}, nil, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE $1 = ANY (tag_ids)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}

func GetPurchase(purchaseId int) (misc.Purchase, error, int) {
	if !misc.IsIdValid(purchaseId) {
		return misc.Purchase{}, errors.New("Not positive id"), misc.NoElement
	}

	p := misc.Purchase{}
	if err := Db.QueryRow(`
		SELECT image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
		if err == sql.ErrNoRows {
			return misc.Purchase{}, err, misc.NoElement
		}

		return misc.Purchase{}, err, misc.NothingToReport
	}

	p.Id = purchaseId
	return p, nil, misc.NothingToReport
}

func LikePurchase(purchaseId, userId int) (error, int) {
	if !misc.IsIdValid(purchaseId) {
		return errors.New("Not positive id"), misc.NoPurchase
	}

	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), misc.VoteForYourself
	}

	// now allow the person to vote for someones else purchase
	sqlResult, err := Db.Exec(`
		INSERT INTO likes (purchase_id, user_id)
		VALUES ($1, $2)`, purchaseId, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num + 1
		WHERE id = $1`, purchaseId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	return nil, misc.NothingToReport
}

func UnlikePurchase(purchaseId, userId int) (error, int) {
	if !misc.IsIdValid(purchaseId) {
		return errors.New("Not positive id"), misc.NoPurchase
	}

	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), misc.VoteForYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM likes
		WHERE purchase_id = $1 AND user_id = $2`, purchaseId, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num - 1
		WHERE id = $1`, purchaseId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	return nil, misc.NothingToReport
}

func AskQuestion(purchaseId, userId int, question string) (int, error, int) {
	// TODO validate everything
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase == userId {
		return 0, errors.New("can't ask question about your stuff"), misc.AskYourself
	}

	questionId := 0
	err = Db.QueryRow(`
		INSERT INTO questions (user_id, purchase_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, purchaseId, question,
	).Scan(&questionId)
	if err != nil {
		err, code := checkSpecificDriverErrors(err)
		return 0, err, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET questions_num = questions_num + 1
		WHERE id = $1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, err, code
	}

	return questionId, nil, misc.NothingToReport
}

// --- Answers ---
func AnswerQuestion(questionId, userId int, answer string) (int, error, int) {
	// TODO validate everything
	whosePurchase, err, code := whoCreatedPurchaseByQuestionId(questionId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase != userId {
		return 0, errors.New("can asnwer only questions regarding your purchase"), misc.AnswerOtherPurchase
	}

	answerId := 0
	err = Db.QueryRow(`
		INSERT INTO answers (user_id, question_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, questionId, answer,
	).Scan(&answerId)
	if err != nil {
		err, code := checkSpecificDriverErrors(err)
		return 0, err, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET answers_num = answers_num + 1
		WHERE id = $1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, err, code
	}

	return answerId, nil, misc.NothingToReport
}
