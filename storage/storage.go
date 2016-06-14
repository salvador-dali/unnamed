// Package storage is responsible for all database operations
package storage

import (
	"../../unnamed/errorCodes"
	"../../unnamed/structs"
	"database/sql"
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
		return id, err, errorCodes.DbNothingToReport
	}

	if errPg, ok := err.(*pq.Error); ok {
		s := string(errPg.Code)
		if s == "23505" {
			return id, err, errorCodes.DbDuplicate
		} else if s == "22001" {
			return id, err, errorCodes.DbValueTooLong
		}
	}

	return id, err, errorCodes.DbNothingToReport
}

func UpdateBrand(id int, name string) (error, int) {
	_, err := Db.Query("UPDATE brands SET name=$1 WHERE id=$2", name, id)
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
		return id, err, errorCodes.DbNothingToReport
	}

	if errPg, ok := err.(*pq.Error); ok {
		s := string(errPg.Code)
		if s == "23505" {
			return id, err, errorCodes.DbDuplicate
		} else if s == "22001" {
			return id, err, errorCodes.DbValueTooLong
		}
	}

	return id, err, errorCodes.DbNothingToReport
}

func UpdateTag(id int, name, descr string) (error, int) {
	_, err := Db.Query("UPDATE tags SET name=$1, description=$2 WHERE id=$3", name, descr, id)
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

// --- Users ---

func GetUser(id int) (structs.User, error, int) {
	user := structs.User{}
	if err := Db.QueryRow("SELECT nickname, image, about, expertise, followers_num, following_num, purchases_num, questions_num, answers_num, issued_at FROM users WHERE id = $1", id).Scan(
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
