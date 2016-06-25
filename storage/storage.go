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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
		log.Println(err)
		return err, misc.NothingToReport
	}

	if affectedRows == 1 {
		return nil, misc.NothingToReport
	} else if affectedRows == 0 {
		log.Println("nothing updated")
		return errors.New("nothing updated"), misc.NothingUpdated
	}
	log.Println(fmt.Sprintf("Expected to update 1 value. %d updated", affectedRows))
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

func GetAllBrands() ([]*misc.Brand, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM brands
		WHERE id > 0`)
	if err != nil {
		log.Println(err)
		return []*misc.Brand{}, misc.NothingToReport
	}
	defer rows.Close()

	brands := []*misc.Brand{}
	for rows.Next() {
		brand := misc.Brand{}
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			log.Println(err)
			return []*misc.Brand{}, misc.NothingToReport
		}
		brands = append(brands, &brand)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.Brand{}, misc.NothingToReport
	}

	return brands, misc.NothingToReport
}

func GetBrand(brandId int) (misc.Brand, int) {
	if !misc.IsIdValid(brandId) {
		log.Println("BrandId is not correct", brandId)
		return misc.Brand{}, misc.NoElement
	}

	brand := misc.Brand{}
	if err := Db.QueryRow(`
		SELECT name, issued_at
		FROM brands
		WHERE id = $1`, brandId,
	).Scan(&brand.Name, &brand.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Brand{}, misc.NoElement
		}

		log.Println(err)
		return misc.Brand{}, misc.NothingToReport
	}

	brand.Id = brandId
	return brand, misc.NothingToReport
}

func CreateBrand(name string) (int, int) {
	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Wrong name for a brand", name)
		return 0, misc.WrongName
	}

	brandId := 0
	err := Db.QueryRow(`
		INSERT INTO brands (name)
		VALUES ($1)
		RETURNING id`, name,
	).Scan(&brandId)
	if err == nil {
		return brandId, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	log.Println("Error creating a brand", err)
	return 0, code
}

func UpdateBrand(brandId int, name string) int {
	if !misc.IsIdValid(brandId) {
		log.Println("BrandId is not correct", brandId)
		return misc.NothingUpdated
	}

	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Brand name is not correct", name)
		return misc.WrongName
	}

	sqlResult, err := Db.Exec(`
		UPDATE brands
		SET name = $1
		WHERE id = $2`, name, brandId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := isAffectedOneRow(sqlResult)
	return code
}

// --- Tags ---

func GetAllTags() ([]*misc.Tag, int) {
	rows, err := Db.Query(`
		SELECT id, name
		FROM tags`)
	if err != nil {
		log.Println(err)
		return []*misc.Tag{}, misc.NothingToReport
	}
	defer rows.Close()

	tags := []*misc.Tag{}
	for rows.Next() {
		tag := misc.Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Println(err)
			return []*misc.Tag{}, misc.NothingToReport
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.Tag{}, misc.NothingToReport
	}

	return tags, misc.NothingToReport
}

func GetTag(tagId int) (misc.Tag, int) {
	if !misc.IsIdValid(tagId) {
		log.Println("TagId is not correct", tagId)
		return misc.Tag{}, misc.NoElement
	}

	tag := misc.Tag{}
	if err := Db.QueryRow(`
		SELECT name, description, issued_at
		FROM tags
		WHERE id = $1`, tagId,
	).Scan(&tag.Name, &tag.Description, &tag.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Tag{}, misc.NoElement
		}

		log.Println(err)
		return misc.Tag{}, misc.NothingToReport
	}

	tag.Id = tagId
	return tag, misc.NothingToReport
}

func CreateTag(name, descr string) (int, int) {
	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Wrong name for a tag", name)
		return 0, misc.WrongName
	}

	descr, ok = misc.ValidateString(descr, misc.MaxLenB)
	if !ok {
		log.Println("Wrong descr for a tag", descr)
		return 0, misc.WrongDescr
	}

	tagId := 0
	err := Db.QueryRow(`
		INSERT INTO tags (name, description)
		VALUES ($1, $2)
		RETURNING id`, name, descr,
	).Scan(&tagId)
	if err == nil {
		return tagId, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	log.Println(err)
	return 0, code
}

func UpdateTag(tagId int, name, descr string) int {
	if !misc.IsIdValid(tagId) {
		log.Println("Tag was not updated", tagId)
		return misc.NothingUpdated
	}

	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Tag has wrong name", name)
		return misc.WrongName
	}

	descr, ok = misc.ValidateString(descr, misc.MaxLenB)
	if !ok {
		log.Println("Tag has wrong description", descr)
		return misc.WrongDescr
	}

	sqlResult, err := Db.Exec(`
		UPDATE tags
		SET name = $1, description = $2
		WHERE id = $3`, name, descr, tagId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := isAffectedOneRow(sqlResult)
	return code
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

func GetUser(userId int) (misc.User, int) {
	if !misc.IsIdValid(userId) {
		log.Println("UserId is not correct")
		return misc.User{}, misc.NoElement
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
			log.Println(err)
			return misc.User{}, misc.NoElement
		}

		log.Println(err)
		return misc.User{}, misc.NothingToReport
	}
	user.Id = userId
	return user, misc.NothingToReport
}

func UpdateUser(userId int, nickname, about string) int {
	if !misc.IsIdValid(userId) {
		log.Println("user was not updated", userId)
		return misc.NothingUpdated
	}

	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		log.Println("Nickname is not correct", nickname)
		return misc.WrongName
	}

	about, ok = misc.ValidateString(about, misc.MaxLenB)
	if !ok {
		log.Println("About is not correct", about)
		return misc.WrongDescr
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET nickname = $1, about = $2
		WHERE id = $3`, nickname, about, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := isAffectedOneRow(sqlResult)
	return code
}

func Follow(whoId, whomId int) int {
	if !misc.IsIdValid(whomId) {
		log.Println("User id is not correct", whomId)
		return misc.NoElement
	}

	if whoId == whomId {
		log.Println("can't follow yourself")
		return misc.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		INSERT INTO followers (who_id, whom_id)
		VALUES ($1, $2)`, whoId, whomId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num + 1
		WHERE id = $1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num + 1
		WHERE id = $1`, whoId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

func Unfollow(whoId, whomId int) int {
	if !misc.IsIdValid(whomId) {
		log.Println("User id is not correct", whomId)
		return misc.NoElement
	}

	if whoId == whomId {
		log.Println("can't follow yourself")
		return misc.FollowYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM followers
		WHERE who_id = $1 AND whom_id = $2`, whoId, whomId)
	if err != nil {
		log.Println(err)
		return misc.NothingToReport
	}

	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num - 1
		WHERE id = $1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num - 1
		WHERE id = $1`, whoId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

func GetFollowing(userId int) ([]*misc.User, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, misc.NothingToReport
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
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			log.Println(err)
			return []*misc.User{}, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}

	return users, misc.NothingToReport
}

func GetFollowers(userId int) ([]*misc.User, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, misc.NothingToReport
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
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			log.Println(err)
			return []*misc.User{}, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}

	return users, misc.NothingToReport
}

func CreateUser(nickname, email, password string) (int, int) {
	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		log.Println("Wrong nickname", nickname)
		return 0, misc.WrongName
	}

	email, ok = misc.ValidateEmail(email)
	if !ok {
		log.Println("Wrong email", email)
		return 0, misc.WrongEmail
	}

	if !misc.IsPasswordValid(password) {
		log.Println("Wrong password")
		return 0, misc.WrongPassword
	}

	salt, err := auth.GenerateSalt()
	if err != nil {
		log.Println(err)
		return 0, misc.NoSalt
	}

	hash, err := auth.PasswordHash(password, salt)
	if err != nil {
		log.Println(err)
		return 0, misc.NothingToReport
	}

	userId := 0
	err = Db.QueryRow(`
		INSERT INTO users (nickname, email, password, salt)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, nickname, email, hash, salt,
	).Scan(&userId)
	if err == nil {
		return userId, misc.NothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	log.Println(err)
	return 0, code
}

func Login(email, password string) (string, bool) {
	email, ok := misc.ValidateEmail(email)
	if !ok || !misc.IsPasswordValid(password) {
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

func getPurchases(rows *sql.Rows, err error) ([]*misc.Purchase, int) {
	if err != nil {
		log.Println(err)
		return []*misc.Purchase{}, misc.NothingToReport
	}
	defer rows.Close()

	purchases := []*misc.Purchase{}
	for rows.Next() {
		p := misc.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			log.Println(err)
			return []*misc.Purchase{}, misc.NothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	return purchases, misc.NothingToReport
}

func whoCreatedPurchaseByPurchaseId(purchaseId int) (int, int) {
	if !misc.IsIdValid(purchaseId) {
		log.Println("purchase ID is wrong", purchaseId)
		return 0, misc.NoPurchase
	}

	whosePurchase := 0
	if err := Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return 0, misc.NoPurchase
		}

		log.Println(err)
		return 0, misc.NothingToReport
	}

	return whosePurchase, misc.NothingToReport
}

func whoCreatedPurchaseByQuestionId(questionId int) (int, int) {
	if !misc.IsIdValid(questionId) {
		// if question does not exist, surely there is no purchase for this question
		log.Println("No question ID is wrong", questionId)
		return 0, misc.NoPurchase
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
			log.Println(err)
			return 0, misc.NoPurchase
		}

		log.Println(err)
		return 0, misc.NothingToReport
	}

	return whosePurchase, misc.NothingToReport
}

func GetUserPurchases(userId int) ([]*misc.Purchase, int) {
	// userId is the current user and is always valid
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE user_id = $1
		ORDER BY issued_at DESC`, userId)

	return getPurchases(rows, err)
}

func CreatePurchase(userId int, description string, brandId int, tagsId []int) (int, int) {
	// userID is the current user and should be valid
	description, ok := misc.ValidateString(description, misc.MaxLenB)
	if !ok {
		log.Println("description is wrong", description)
		return 0, misc.WrongDescr
	}

	if brandId < 0 {
		log.Println("BrandID is wrong", brandId)
		return 0, misc.NoElement
	}

	if err, code := validateTags(tagsId); err != nil {
		log.Println(err)
		return 0, code
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
		log.Println(err)
		return 0, misc.NothingToReport
	}

	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return 0, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET purchases_num = purchases_num + 1
		WHERE id=$1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return id, misc.NothingToReport
}

func GetAllPurchases() ([]*misc.Purchase, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithBrand(brandId int) ([]*misc.Purchase, int) {
	if !misc.IsIdValid(brandId) {
		log.Println("Brand Id is wrong", brandId)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE brand_id = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithTag(tagId int) ([]*misc.Purchase, int) {
	if !misc.IsIdValid(tagId) {
		log.Println("Tag ID is wrong", tagId)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE $1 = ANY (tag_ids)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}

func GetPurchase(purchaseId int) (misc.Purchase, int) {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase ID is wrong", purchaseId)
		return misc.Purchase{}, misc.NoElement
	}

	p := misc.Purchase{}
	if err := Db.QueryRow(`
		SELECT image, description, user_id, issued_at, brand_id, likes_num
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Purchase{}, misc.NoElement
		}

		log.Println(err)
		return misc.Purchase{}, misc.NothingToReport
	}

	p.Id = purchaseId
	return p, misc.NothingToReport
}

func LikePurchase(purchaseId, userId int) int {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase Id is not valid", purchaseId)
		return misc.NoPurchase
	}

	// check whose purchase is it
	whosePurchase, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return misc.VoteForYourself
	}

	// now allow the person to vote for someones else purchase
	sqlResult, err := Db.Exec(`
		INSERT INTO likes (purchase_id, user_id)
		VALUES ($1, $2)`, purchaseId, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num + 1
		WHERE id = $1`, purchaseId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

func UnlikePurchase(purchaseId, userId int) int {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase Id is not possitive", purchaseId)
		return misc.NoPurchase
	}

	// check whose purchase is it
	whosePurchase, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return misc.VoteForYourself
	}

	sqlResult, err := Db.Exec(`
		DELETE FROM likes
		WHERE purchase_id = $1 AND user_id = $2`, purchaseId, userId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num - 1
		WHERE id = $1`, purchaseId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

func AskQuestion(purchaseId, userId int, question string) (int, int) {
	// TODO validate everything
	whosePurchase, code := whoCreatedPurchaseByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return 0, code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return 0, misc.AskYourself
	}

	questionId := 0
	err := Db.QueryRow(`
		INSERT INTO questions (user_id, purchase_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, purchaseId, question,
	).Scan(&questionId)
	if err != nil {
		err, code := checkSpecificDriverErrors(err)
		log.Println(err)
		return 0, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET questions_num = questions_num + 1
		WHERE id = $1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return questionId, misc.NothingToReport
}

// --- Answers ---
func AnswerQuestion(questionId, userId int, answer string) (int, int) {
	// TODO validate everything
	whosePurchase, code := whoCreatedPurchaseByQuestionId(questionId)
	if whosePurchase == 0 {
		return 0, code
	}

	if whosePurchase != userId {
		log.Println("can asnwer only questions regarding your purchase")
		return 0, misc.AnswerOtherPurchase
	}

	answerId := 0
	err := Db.QueryRow(`
		INSERT INTO answers (user_id, question_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, questionId, answer,
	).Scan(&answerId)
	if err != nil {
		err, code := checkSpecificDriverErrors(err)
		log.Println(err)
		return 0, code
	}

	sqlResult, err := Db.Exec(`
		UPDATE users
		SET answers_num = answers_num + 1
		WHERE id = $1`, userId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return answerId, misc.NothingToReport
}
