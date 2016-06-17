// Package storage is responsible for all database operations
package storage

import (
	"../../unnamed/errorCodes"
	"../../unnamed/structs"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
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
		return err, errorCodes.DbNothingToReport
	}

	if affectedRows == 1 {
		return nil, errorCodes.DbNothingToReport
	} else if affectedRows == 0 {
		return errors.New("nothing updated"), errorCodes.DbNothingUpdated
	}
	return errors.New(fmt.Sprintf("Expected to update 1 value. %d updated", affectedRows)), errorCodes.DbNothingToReport
}

// checkSpecificDriverErrors analyses the error result against specific errors that a client should know about
// This checks for Value Limit violation and Duplicate constraint violation
func checkSpecificDriverErrors(err error) (error, int) {
	if errPg, ok := err.(*pq.Error); ok {
		s := string(errPg.Code)
		switch {
		case s == "23505":
			return err, errorCodes.DbDuplicate
		case s == "22001":
			return err, errorCodes.DbValueTooLong
		}
	}

	return err, errorCodes.DbNothingToReport
}

// --- Brands ---

func GetAllBrands() ([]*structs.Brand, error, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM brands
		WHERE id > 0`)
	if err != nil {
		return []*structs.Brand{}, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	brands := []*structs.Brand{}
	for rows.Next() {
		brand := structs.Brand{}
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			return []*structs.Brand{}, err, errorCodes.DbNothingToReport
		}
		brands = append(brands, &brand)
	}

	if err = rows.Err(); err != nil {
		return []*structs.Brand{}, err, errorCodes.DbNothingToReport
	}

	return brands, nil, errorCodes.DbNothingToReport
}

func GetBrand(brandId int) (structs.Brand, error, int) {
	if brandId <= 0 {
		return structs.Brand{}, errors.New("Not positive id"), errorCodes.DbNoElement
	}

	brand := structs.Brand{}
	if err := Db.QueryRow(`
		SELECT name, issued_at
		FROM brands
		WHERE id = $1`, brandId,
	).Scan(&brand.Name, &brand.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return structs.Brand{}, err, errorCodes.DbNoElement
		}

		return structs.Brand{}, err, errorCodes.DbNothingToReport
	}

	brand.Id = brandId
	return brand, nil, errorCodes.DbNothingToReport
}

func CreateBrand(name string) (int, error, int) {
	brandId := 0
	err := Db.QueryRow(`
		INSERT INTO brands (name)
		VALUES ($1)
		RETURNING id`, name,
	).Scan(&brandId)
	if err == nil {
		return brandId, nil, errorCodes.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateBrand(brandId int, name string) (error, int) {
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

func GetAllTags() ([]*structs.Tag, error, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM tags`)
	if err != nil {
		return []*structs.Tag{}, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	tags := []*structs.Tag{}
	for rows.Next() {
		tag := structs.Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return []*structs.Tag{}, err, errorCodes.DbNothingToReport
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return []*structs.Tag{}, err, errorCodes.DbNothingToReport
	}

	return tags, nil, errorCodes.DbNothingToReport
}

func GetTag(tagId int) (structs.Tag, error, int) {
	tag := structs.Tag{}
	if err := Db.QueryRow(`
		SELECT name, description, issued_at
		FROM tags
		WHERE id = $1`, tagId,
	).Scan(&tag.Name, &tag.Description, &tag.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return structs.Tag{}, err, errorCodes.DbNoElement
		}

		return structs.Tag{}, err, errorCodes.DbNothingToReport
	}

	tag.Id = tagId
	return tag, nil, errorCodes.DbNothingToReport
}

func CreateTag(name, descr string) (int, error, int) {
	tagId := 0
	err := Db.QueryRow(`
		INSERT INTO tags (name, description)
		VALUES ($1, $2)
		RETURNING id`, name, descr,
	).Scan(&tagId)
	if err == nil {
		return tagId, nil, errorCodes.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return 0, err, code
}

func UpdateTag(tagId int, name, descr string) (error, int) {
	sqlResult, err := Db.Exec(`
		UPDATE tags
		SET name=$1, description=$2
		WHERE id=$3`, name, descr, tagId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

// validateTags makes sure that all the tags are in the database. psql does not support this http://dba.stackexchange.com/q/60132/15318
func validateTags(tagIds []int) (error, int) {
	if len(tagIds) == 0 {
		return errors.New("no tags"), errorCodes.NoTags
	}
	if len(tagIds) > maxTags {
		return errors.New("too many tags"), errorCodes.TooManyTags
	}

	for _, v := range tagIds {
		if v <= 0 {
			return errors.New("tag is negative"), errorCodes.IdNotNatural
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
		return err, errorCodes.DbNothingToReport
	}

	if num != len(tagIds) {
		return errors.New("some tags are missing"), errorCodes.DbNotAllTagsCorrect
	}

	return nil, errorCodes.DbNothingToReport
}

// --- Users ---

func GetUser(userId int) (structs.User, error, int) {
	user := structs.User{}
	if err := Db.QueryRow(`
		SELECT nickname, image, about, expertise, followers_num, following_num, purchases_num, questions_num, answers_num, issued_at
		FROM users WHERE id = $1`, userId,
	).Scan(
		&user.Nickname, &user.Image, &user.About, &user.Expertise, &user.Followers_num,
		&user.Following_num, &user.Purchases_num, &user.Questions_num, &user.Answers_num,
		&user.Issued_at,
	); err != nil {
		if err == sql.ErrNoRows {
			return structs.User{}, err, errorCodes.DbNoElement
		}

		return structs.User{}, err, errorCodes.DbNothingToReport
	}
	user.Id = userId
	return user, nil, errorCodes.DbNothingToReport
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
		return errors.New("can't follow yourself"), errorCodes.FollowYourself
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

	return nil, errorCodes.DbNothingToReport
}

func Unfollow(whoId, whomId int) (error, int) {
	if whoId == whomId {
		return errors.New("can't follow yourself"), errorCodes.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM followers
		WHERE who_id = $1 AND whom_id = $2`, whoId, whomId)
	if err != nil {
		return err, errorCodes.DbNothingToReport
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

	return nil, errorCodes.DbNothingToReport
}

func GetFollowing(userId int) ([]*structs.User, error, int) {
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT whom_id
			FROM followers
			WHERE who_id = $1
		)`, userId)
	if err != nil {
		return []*structs.User{}, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	users := []*structs.User{}
	for rows.Next() {
		user := structs.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*structs.User{}, err, errorCodes.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*structs.User{}, err, errorCodes.DbNothingToReport
	}

	return users, nil, errorCodes.DbNothingToReport
}

func GetFollowers(userId int) ([]*structs.User, error, int) {
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT who_id
			FROM followers
			WHERE whom_id = $1
		)`, userId)
	if err != nil {
		return []*structs.User{}, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	users := []*structs.User{}
	for rows.Next() {
		user := structs.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return []*structs.User{}, err, errorCodes.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return []*structs.User{}, err, errorCodes.DbNothingToReport
	}

	return users, nil, errorCodes.DbNothingToReport
}

// --- Purchases ---

func getPurchases(rows *sql.Rows, err error) ([]*structs.Purchase, error, int) {
	if err != nil {
		return []*structs.Purchase{}, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	purchases := []*structs.Purchase{}
	for rows.Next() {
		p := structs.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			return []*structs.Purchase{}, err, errorCodes.DbNothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return []*structs.Purchase{}, err, errorCodes.DbNothingToReport
	}

	return purchases, nil, errorCodes.DbNothingToReport
}

func GetUserPurchases(userId int) ([]*structs.Purchase, error, int) {
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
		return 0, err, errorCodes.DbNothingToReport
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

	return id, nil, errorCodes.DbNothingToReport
}

func GetAllPurchases() ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithBrand(brandId int) ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE brand = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithTag(tagId int) ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE $1 = ANY (tags)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}

func GetPurchase(purchaseId int) (structs.Purchase, error, int) {
	p := structs.Purchase{}
	if err := Db.QueryRow(`
		SELECT image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
		if err == sql.ErrNoRows {
			return structs.Purchase{}, err, errorCodes.DbNoElement
		}

		return structs.Purchase{}, err, errorCodes.DbNothingToReport
	}

	p.Id = purchaseId
	return p, nil, errorCodes.DbNothingToReport
}

func whoCreatedPurchaseByPurchaseId(purchaseId int) (int, error, int) {
	whosePurchase := 0
	if err := Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			return 0, err, errorCodes.DbNoPurchase
		}

		return 0, err, errorCodes.DbNothingToReport
	}

	return whosePurchase, nil, errorCodes.DbNothingToReport
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
			return 0, err, errorCodes.DbNoPurchaseForQuestion
		}

		return 0, err, errorCodes.DbNothingToReport
	}

	return whosePurchase, nil, errorCodes.DbNothingToReport
}

func LikePurchase(purchaseId, userId int) (error, int) {
	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), errorCodes.DbVoteForOwnStuff
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

	return nil, errorCodes.DbNothingToReport
}

func UnlikePurchase(purchaseId, userId int) (error, int) {
	// check whose purchase it is
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return err, code
	}

	if whosePurchase == userId {
		return errors.New("can't vote for own purchase"), errorCodes.DbVoteForOwnStuff
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

	return nil, errorCodes.DbNothingToReport
}

func AskQuestion(purchaseId, userId int, question string) (int, error, int) {
	whosePurchase, err, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase == userId {
		return 0, errors.New("can't ask question about your stuff"), errorCodes.DbAskAboutOwnStuff
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

	return questionId, nil, errorCodes.DbNothingToReport
}

// --- Answers ---
func AnswerQuestion(questionId, userId int, answer string) (int, error, int) {
	whosePurchase, err, code := whoCreatedPurchaseByQuestionId(questionId)
	if err != nil {
		return 0, err, code
	}

	if whosePurchase != userId {
		return 0, errors.New("can asnwer only questions regarding your purchase"), errorCodes.DbCannotAnswerOtherPurchase
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

	return answerId, nil, errorCodes.DbNothingToReport
}
