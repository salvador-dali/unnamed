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
		if s == "23505" {
			return err, errorCodes.DbDuplicate
		} else if s == "22001" {
			return err, errorCodes.DbValueTooLong
		}
	}

	return err, errorCodes.DbNothingToReport
}

// --- Brands ---

func GetAllBrands() ([]*structs.Brand, error, int) {
	brands := []*structs.Brand{}

	rows, err := Db.Query("SELECT id, name FROM brands")
	if err != nil {
		return brands, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		brand := structs.Brand{}
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			return brands, err, errorCodes.DbNothingToReport
		}
		brands = append(brands, &brand)
	}

	if err = rows.Err(); err != nil {
		return brands, err, errorCodes.DbNothingToReport
	}

	return brands, nil, errorCodes.DbNothingToReport
}

func GetBrand(id int) (structs.Brand, error, int) {
	brand := structs.Brand{}
	if err := Db.QueryRow("SELECT name, issued_at FROM brands WHERE id = $1", id).Scan(&brand.Name, &brand.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return brand, err, errorCodes.DbNoElement
		}

		return brand, err, errorCodes.DbNothingToReport
	}

	brand.Id = id
	return brand, nil, errorCodes.DbNothingToReport
}

func CreateBrand(name string) (int, error, int) {
	id := 0
	err := Db.QueryRow("INSERT INTO brands (name) VALUES($1) RETURNING id", name).Scan(&id)
	if err == nil {
		return id, nil, errorCodes.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return id, err, code
}

func UpdateBrand(id int, name string) (error, int) {
	sqlResult, err := Db.Exec("UPDATE brands SET name=$1 WHERE id=$2", name, id)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

// --- Tags ---

func GetAllTags() ([]*structs.Tag, error, int) {
	tags := []*structs.Tag{}

	rows, err := Db.Query("SELECT id, name FROM tags")
	if err != nil {
		return tags, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		tag := structs.Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return tags, err, errorCodes.DbNothingToReport
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return tags, err, errorCodes.DbNothingToReport
	}

	return tags, nil, errorCodes.DbNothingToReport
}

func GetTag(id int) (structs.Tag, error, int) {
	tag := structs.Tag{}
	if err := Db.QueryRow("SELECT name, description, issued_at FROM tags WHERE id = $1", id).Scan(&tag.Name, &tag.Description, &tag.Issued_at); err != nil {
		if err == sql.ErrNoRows {
			return tag, err, errorCodes.DbNoElement
		}

		return tag, err, errorCodes.DbNothingToReport
	}

	tag.Id = id
	return tag, nil, errorCodes.DbNothingToReport
}

func CreateTag(name, descr string) (int, error, int) {
	id := 0
	err := Db.QueryRow("INSERT INTO tags (name, description) VALUES($1, $2) RETURNING id", name, descr).Scan(&id)
	if err == nil {
		return id, nil, errorCodes.DbNothingToReport
	}

	err, code := checkSpecificDriverErrors(err)
	return id, err, code
}

func UpdateTag(id int, name, descr string) (error, int) {
	sqlResult, err := Db.Exec(`
		UPDATE tags
		SET name=$1, description=$2
		WHERE id=$3`, name, descr, id)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}

	return isAffectedOneRow(sqlResult)
}

// --- Users ---

func GetUser(id int) (structs.User, error, int) {
	user := structs.User{}
	if err := Db.QueryRow(`
		SELECT nickname, image, about, expertise, followers_num, following_num, purchases_num, questions_num, answers_num, issued_at
		FROM users WHERE id = $1`, id).Scan(
		&user.Nickname, &user.Image, &user.About, &user.Expertise, &user.Followers_num, &user.Following_num,
		&user.Purchases_num, &user.Questions_num, &user.Answers_num, &user.Issued_at,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, err, errorCodes.DbNoElement
		}

		return user, err, errorCodes.DbNothingToReport
	}
	user.Id = id
	return user, nil, errorCodes.DbNothingToReport
}

func UpdateUser(id int, nickname, about string) (error, int) {
	res, err := Db.Exec(`
		UPDATE users
		SET nickname=$1, about=$2
		WHERE id=$3", nickname, about, id`)
	if errPg, ok := err.(*pq.Error); ok {
		s := string(errPg.Code)
		if s == "23505" {
			return err, errorCodes.DbDuplicate
		} else if s == "22001" {
			return err, errorCodes.DbValueTooLong
		}
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err, errorCodes.DbNothingToReport
	}

	if affect == 0 {
		return errors.New("nothing updated"), errorCodes.DbNothingUpdated
	}
	return nil, errorCodes.DbNothingToReport
}

func Follow(whoId, whomId int) (error, int) {
	if whoId == whomId {
		return errors.New("can't follow yourself"), errorCodes.FollowYourself
	}

	sqlResult, err := Db.Exec("INSERT INTO followers (who_id, whom_id) VALUES($1, $2)", whoId, whomId)
	if err, code := checkSpecificDriverErrors(err); err != nil {
		return err, code
	}
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num + 1
		WHERE id=$1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num + 1
		WHERE id=$1`, whoId)
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
		WHERE who_id=$1 AND whom_id=$2`, whoId, whomId)
	if err != nil {
		return err, errorCodes.DbNothingToReport
	}

	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET followers_num = followers_num - 1
		WHERE id=$1`, whomId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	sqlResult, err = Db.Exec(`
		UPDATE users
		SET following_num = following_num - 1
		WHERE id=$1`, whoId)
	if err, code := isAffectedOneRow(sqlResult); err != nil {
		return err, code
	}

	return nil, errorCodes.DbNothingToReport
}

func GetFollowering(id int) ([]*structs.User, error, int) {
	users := []*structs.User{}
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT whom_id
			FROM followers
			WHERE who_id = $1
		)`, id)
	if err != nil {
		return users, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		user := structs.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return users, err, errorCodes.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return users, err, errorCodes.DbNothingToReport
	}

	return users, nil, errorCodes.DbNothingToReport
}

func GetFollowers(id int) ([]*structs.User, error, int) {
	users := []*structs.User{}
	rows, err := Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT who_id
			FROM followers
			WHERE whom_id = $1
		)`, id)
	if err != nil {
		return users, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		user := structs.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			return users, err, errorCodes.DbNothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return users, err, errorCodes.DbNothingToReport
	}

	return users, nil, errorCodes.DbNothingToReport
}

func GetUserPurchases(id int) ([]*structs.Purchase, error, int) {
	purchases := []*structs.Purchase{}
	rows, err := Db.Query(`
		SELECT id, image, description, issued_at, brand, likes_num
		FROM purchases
		WHERE user_id = $1
		ORDER BY issued_at DESC`, id)
	if err != nil {
		return purchases, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		p := structs.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			return purchases, err, errorCodes.DbNothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return purchases, err, errorCodes.DbNothingToReport
	}

	return purchases, nil, errorCodes.DbNothingToReport
}

// --- Purchases ---

func getPurchases(rows *sql.Rows, err error) ([]*structs.Purchase, error, int) {
	purchases := []*structs.Purchase{}
	if err != nil {
		return purchases, err, errorCodes.DbNothingToReport
	}
	defer rows.Close()

	for rows.Next() {
		p := structs.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &p.Issued_at, &p.Brand, &p.Likes_num); err != nil {
			return purchases, err, errorCodes.DbNothingToReport
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return purchases, err, errorCodes.DbNothingToReport
	}

	return purchases, nil, errorCodes.DbNothingToReport
}

func GetAllPurchases() ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithBrand(brandId int) ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand, likes_num
		FROM purchases
		WHERE brand = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

func GetAllPurchasesWithTag(tagId int) ([]*structs.Purchase, error, int) {
	rows, err := Db.Query(`
		SELECT id, image, description, user_id, issued_at, brand, likes_num
		FROM purchases
		WHERE $1 = ANY (tags)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}
