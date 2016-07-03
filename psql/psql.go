// Package storage is responsible for all database operations
package psql

import (
	"../config"
	"../misc"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
)

var Db *sql.DB

// Init prepares the database abstraction for later use
func Init() {
	// It does not establish any connections to the database, nor does it validate driver
	// connection parameters. To do this call Ping http://go-database-sql.org/accessing.html
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Cfg.DbUser,
		config.Cfg.DbPass,
		config.Cfg.DbHost,
		config.Cfg.DbPort,
		config.Cfg.DbName,
	)

	if db, err := sql.Open("postgres", dbURL); err != nil {
		log.Fatal(err)
	} else {
		Db = db
	}
}

// IsAffectedOneRow checks that the result of a query executed with Exec has modified only 1 row
func IsAffectedOneRow(sqlResult sql.Result) (error, int) {
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

// CheckSpecificDriverErrors analyses the error result against specific errors that a client should know about
// This checks for Value Limit violation and Duplicate constraint violation
func CheckSpecificDriverErrors(err error) (error, int) {
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
