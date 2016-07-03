package tag

import (
	"../../misc"
	"../../psql"
	"bytes"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"
)

// ShowAll returns a list of all possible tags
func ShowAll() ([]*misc.Tag, int) {
	rows, err := psql.Db.Query(`
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

// Show a tag by Id
func ShowById(tagId int) (misc.Tag, int) {
	if !misc.IsIdValid(tagId) {
		log.Println("TagId is not correct", tagId)
		return misc.Tag{}, misc.NoElement
	}

	tag := misc.Tag{}
	var timestamp time.Time
	if err := psql.Db.QueryRow(`
		SELECT name, description, issued_at
		FROM tags
		WHERE id = $1`, tagId,
	).Scan(&tag.Name, &tag.Description, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Tag{}, misc.NoElement
		}

		log.Println(err)
		return misc.Tag{}, misc.NothingToReport
	}

	tag.Id = tagId
	tag.Issued_at = timestamp.Unix()
	return tag, misc.NothingToReport
}

// Create a new tag
func Create(name, descr string) (int, int) {
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
	err := psql.Db.QueryRow(`
		INSERT INTO tags (name, description)
		VALUES ($1, $2)
		RETURNING id`, name, descr,
	).Scan(&tagId)
	if err == nil {
		return tagId, misc.NothingToReport
	}

	err, code := psql.CheckSpecificDriverErrors(err)
	log.Println(err)
	return 0, code
}

// Update a tag by Id
func Update(tagId int, name, descr string) int {
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

	sqlResult, err := psql.Db.Exec(`
		UPDATE tags
		SET name = $1, description = $2
		WHERE id = $3`, name, descr, tagId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := psql.IsAffectedOneRow(sqlResult)
	return code
}

// ValidateTags makes sure that all the tagIds exist in the database
func ValidateTags(tagIds []int) (error, int) {
	// psql does not support this http://dba.stackexchange.com/q/60132/15318
	if len(tagIds) == 0 {
		return errors.New("no tags"), misc.NoTags
	}

	if len(tagIds) > misc.MaxTags {
		return errors.New("too many tags"), misc.WrongTagsNum
	}

	buf := bytes.NewBufferString("SELECT COUNT(id) FROM tags WHERE id IN (")
	for i, v := range tagIds {
		if v <= 0 {
			return errors.New("tag is negative"), misc.WrongTags
		}
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	num := 0
	if err := psql.Db.QueryRow(buf.String()).Scan(&num); err != nil {
		return err, misc.NothingToReport
	}

	if num != len(tagIds) {
		return errors.New("some tags are missing"), misc.WrongTags
	}

	return nil, misc.NothingToReport
}
