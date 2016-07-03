package brand

import (
	"../../misc"
	"../../psql"
	"database/sql"
	"log"
	"time"
)

// ShowAll returns a list of all possible brands
func ShowAll() ([]*misc.Brand, int) {
	rows, err := psql.Db.Query(`
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

// Show a brand by Id
func ShowById(brandId int) (misc.Brand, int) {
	if !misc.IsIdValid(brandId) {
		log.Println("BrandId is not correct", brandId)
		return misc.Brand{}, misc.NoElement
	}

	brand := misc.Brand{}
	var timestamp time.Time
	if err := psql.Db.QueryRow(`
		SELECT name, issued_at
		FROM brands
		WHERE id = $1`, brandId,
	).Scan(&brand.Name, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Brand{}, misc.NoElement
		}

		log.Println(err)
		return misc.Brand{}, misc.NothingToReport
	}

	brand.Id = brandId
	brand.Issued_at = timestamp.Unix()
	return brand, misc.NothingToReport
}

// Create a new brand
func Create(name string) (int, int) {
	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Wrong name for a brand", name)
		return 0, misc.WrongName
	}

	brandId := 0
	err := psql.Db.QueryRow(`
		INSERT INTO brands (name)
		VALUES ($1)
		RETURNING id`, name,
	).Scan(&brandId)
	if err == nil {
		return brandId, misc.NothingToReport
	}

	err, code := psql.CheckSpecificDriverErrors(err)
	log.Println("Error creating a brand", err)
	return 0, code
}

// Update a brand by Id
func Update(brandId int, name string) int {
	if !misc.IsIdValid(brandId) {
		log.Println("BrandId is not correct", brandId)
		return misc.NothingUpdated
	}

	name, ok := misc.ValidateString(name, misc.MaxLenS)
	if !ok {
		log.Println("Brand name is not correct", name)
		return misc.WrongName
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE brands
		SET name = $1
		WHERE id = $2`, name, brandId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := psql.IsAffectedOneRow(sqlResult)
	return code
}
