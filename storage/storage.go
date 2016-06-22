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

const maxTags = 4 // maximum number of tags for a purchase

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
		return err, misc.DbNothingToReport
	}

	if affectedRows == 1 {
		return nil, misc.DbNothingToReport
	} else if affectedRows == 0 {
		return errors.New("nothing updated"), misc.DbNothingUpdated
	}
	return errors.New(fmt.Sprintf("Expected to update 1 value. %d updated", affectedRows)), misc.DbNothingToReport
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
			return err, misc.DbValueTooLong
		case s == "23503":
			return err, misc.DbForeignKeyViolation
		}
	}

	return err, misc.DbNothingToReport
}

// --- Brands ---

func GetAllBrands() ([]*misc.Brand, error, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM brands
		WHERE id > 0`)
	if err != nil {
		return []*misc.Brand{}, err, misc.DbNothingToReport
	}
	defer rows.Close()

	brands := []*misc.Brand{}
	for rows.Next() {
		brand := misc.Brand{}
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			return []*misc.Brand{}, err, misc.DbNothingToReport
		}
		brands = append(brands, &brand)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Brand{}, err, misc.DbNothingToReport
	}

	return brands, nil, misc.DbNothingToReport
}

func GetBrand(brandId int) (misc.Brand, error, int) {
	if brandId <= 0 {
		return misc.Brand{}, errors.New("Not positive id"), misc.DbNoElement
	}

	brand := misc.Brand{}
	if err := Db.QueryRow(`
		SELECT name, issued_at
		FROM brands
		WHERE id = $1`, brandId,
	).Scan(&brand.Name, &brand.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return misc.Brand{}, err, misc.DbNoElement
		}

		return misc.Brand{}, err, misc.DbNothingToReport
	}

	brand.Id = brandId
	return brand, nil, misc.DbNothingToReport
}

func CreateBrand(name string) (int, error, int) {
	brandId := 0
	err := Db.QueryRow(`
		INSERT INTO brands (name)
		VALUES ($1)
		RETURNING id`, name,
	).Scan(&brandId)
	if err == nil {
		return brandId, nil, misc.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateBrand(brandId int, name string) (error, int) {
	if brandId <= 0 {
		return errors.New("nothing updated"), misc.DbNothingUpdated
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
		return []*misc.Tag{}, err, misc.DbNothingToReport
	}
	defer rows.Close()

	tags := []*misc.Tag{}
	for rows.Next() {
		tag := misc.Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return []*misc.Tag{}, err, misc.DbNothingToReport
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Tag{}, err, misc.DbNothingToReport
	}

	return tags, nil, misc.DbNothingToReport
}

func GetTag(tagId int) (misc.Tag, error, int) {
	tag := misc.Tag{}
	if err := Db.QueryRow(`
		SELECT name, description, issued_at
		FROM tags
		WHERE id = $1`, tagId,
	).Scan(&tag.Name, &tag.Description, &tag.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return misc.Tag{}, err, misc.DbNoElement
		}

		return misc.Tag{}, err, misc.DbNothingToReport
	}

	tag.Id = tagId
	return tag, nil, misc.DbNothingToReport
}

func CreateTag(name, descr string) (int, error, int) {
	tagId := 0
	err := Db.QueryRow(`
		INSERT INTO tags (name, description)
		VALUES ($1, $2)
		RETURNING id`, name, descr,
	).Scan(&tagId)
	if err == nil {
		return tagId, nil, misc.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateTag(tagId int, name, descr string) (error, int) {
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
	if len(tagIds) > maxTags {
		return errors.New("too many tags"), misc.TooManyTags
	}

	for _, v := range tagIds {
		if v <= 0 {
			return errors.New("tag is negative"), misc.IdNotNatural
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
		return err, misc.DbNothingToReport
	}

	if num != len(tagIds) {
		return errors.New("some tags are missing"), misc.DbNotAllTagsCorrect
	}

	return nil, misc.DbNothingToReport
}

// --- Users ---

func GetUser(userId int) (misc.User, error, int) {
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
			return misc.User{}, err, misc.DbNoElement
		}

		return misc.User{}, err, misc.DbNothingToReport
	}
	user.Id = userId
	return user, nil, misc.DbNothingToReport
}

func UpdateUser(userId int, nickname, about string) (error, int) {
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

	return nil, misc.DbNothingToReport
}

func Unfollow(whoId, whomId int) (error, int) {
	if whoId == whomId {
		return errors.New("can't follow yourself"), misc.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM followers
		WHERE who_id = $1 AND whom_id = $2`, whoId, whomId)
	if err != nil {
		return err, misc.DbNothingToReport
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

	return nil, misc.DbNothingToReport
}

func GetFollowing(userId int) ([]*misc.User, error, int) {
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT whom_id
			FROM followers
			WHERE who_id = $1
		)`, userId)
	if err != nil {
		return []*misc.User{}, err, misc.DbNothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*misc.User{}, err, misc.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*misc.User{}, err, misc.DbNothingToReport
	}

	return users, nil, misc.DbNothingToReport
}

func GetFollowers(userId int) ([]*misc.User, error, int) {
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT who_id
			FROM followers
			WHERE whom_id = $1
		)`, userId)
	if err != nil {
		return []*misc.User{}, err, misc.DbNothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*misc.User{}, err, misc.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*misc.User{}, err, misc.DbNothingToReport
	}

	return users, nil, misc.DbNothingToReport
}

func CreateUser(nickname, email, password string) (int, error, int) {
	salt, err := auth.GenerateSalt()
	if err != nil {
		return 0, err, misc.NoSalt
	}

	hash, err := auth.PasswordHash(password, salt)
	if err != nil {
		return 0, err, misc.DbNothingToReport
	}

	userId := 0
	err = Db.QueryRow(`
		INSERT INTO users (nickname, email, password, salt)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, nickname, email, hash, salt,
	).Scan(&userId)
	if err == nil {
		return userId, nil, misc.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func Login(email, password string) (string, bool) {
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
		return []*misc.Purchase{}, err, misc.DbNothingToReport
	}
	defer rows.Close()

	purchases := []*misc.Purchase{}
	for rows.Next() {
		p := misc.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			return []*misc.Purchase{}, err, misc.DbNothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return []*misc.Purchase{}, err, misc.DbNothingToReport
	}

	return purchases, nil, misc.DbNothingToReport
}

func whoCreatedPurchaseByPurchaseId(purchaseId int) (int, error, int) {
	whosePurchase := 0
	if err := Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			return 0, err, misc.DbNoPurchase
		}

		return 0, err, misc.DbNothingToReport
	}

	return whosePurchase, nil, misc.DbNothingToReport
}

func whoCreatedPurchaseByQuestionId(questionId int) (int, error, int) {
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
			return 0, err, misc.DbNoPurchaseForQuestion
		}

		return 0, err, misc.DbNothingToReport
	}

	return whosePurchase, nil, misc.DbNothingToReport
}

func GetUserPurchases(userId int) ([]*misc.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE user_id = $1
		ORDER BY issued_at DESC`, userId)

	return getPurchases(rows, err)
}

func CreatePurchase(userId int, description string, brandId int, tagsId []int) (int, error, int) {
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
		return 0, err, misc.DbNothingToReport
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

	return id, nil, misc.DbNothingToReport
}

func GetAllPurchases() ([]*misc.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithBrand(brandId int) ([]*misc.Purchase, error, int) {
	if brandId <= 0 {
		return []*misc.Purchase{}, nil, misc.DbNothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE brand_id = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithTag(tagId int) ([]*misc.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE $1 = ANY (tag_ids)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}

func GetPurchase(purchaseId int) (misc.Purchase, error, int) {
	p := misc.Purchase{}
	if err := Db.QueryRow(`
		SELECT image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
		if err == sql.ErrNoRows {
			return misc.Purchase{}, err, misc.DbNoElement
		}

		return misc.Purchase{}, err, misc.DbNothingToReport
	}

	p.Id = purchaseId
	return p, nil, misc.DbNothingToReport
}

func LikePurchase(purchaseId, userId int) (error, int) {
	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), misc.DbVoteForOwnStuff
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

	return nil, misc.DbNothingToReport
}

func UnlikePurchase(purchaseId, userId int) (error, int) {
	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), misc.DbVoteForOwnStuff
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

	return nil, misc.DbNothingToReport
}

func AskQuestion(purchaseId, userId int, question string) (int, error, int) {
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase == userId {
		return 0, errors.New("can't ask question about your stuff"), misc.DbAskAboutOwnStuff
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

	return questionId, nil, misc.DbNothingToReport
}

// --- Answers ---
func AnswerQuestion(questionId, userId int, answer string) (int, error, int) {
	whosePurchase, err, code := whoCreatedPurchaseByQuestionId(questionId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase != userId {
		return 0, errors.New("can asnwer only questions regarding your purchase"), misc.DbCannotAnswerOtherPurchase
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

	return answerId, nil, misc.DbNothingToReport
}
