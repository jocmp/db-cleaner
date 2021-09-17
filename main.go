package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Masterminds/squirrel"
	driver "github.com/denisenkom/go-mssqldb"
	sqlx "github.com/jmoiron/sqlx"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

func main() {
	db := MustOpenDatabaseConnection()
	CleanDatabase(db)
}

func CleanDatabase(db *sqlx.DB) {
	databaseTables := mustFindDatabaseTables(db)
	for _, table := range databaseTables {
		mustDeleteFromTable(db, table)
	}
}

func mustFindDatabaseTables(db *sqlx.DB) []DatabaseTable {
	selectTableName, _, err := squirrel.Select("table_name").From("information_schema.tables").ToSql()
	if err != nil {
		log.Fatal(err)
	}
	tables := []DatabaseTable{}
	err = db.Select(&tables, selectTableName)
	if err != nil {
		log.Fatal(err)
	}
	return tables
}

func mustDeleteFromTable(db *sqlx.DB, table DatabaseTable) {
	deleteFromTable, _, err := squirrel.Delete(table.Name).ToSql()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(deleteFromTable)
	if err != nil {
		log.Fatal(err)
	}
}

type DatabaseTable struct {
	Name string `db:"table_name"`
}

func MustOpenDatabaseConnection() *sqlx.DB {
	query := url.Values{}
	query.Add("database", "<YOUR_DATABASE_NAME>")
	connectionURL := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword("<YOUR_USERNAME>", "<YOUR_PASSWORD>"),
		Host:     fmt.Sprintf("%s:%d", "0.0.0.0", 1433),
		RawQuery: query.Encode(),
	}
	log.Printf("Your URL: %s\n", connectionURL.String())
	sqltrace.Register("sqlserver", &driver.Driver{})
	db, err := sqlxtrace.Open("sqlserver", connectionURL.String())
	db.MapperFunc(strings.Title)
	if err != nil {
		log.Fatal(err)
	}
	printDatabaseName(db)
	return db
}

func printDatabaseName(db *sqlx.DB) {
	rows, err := db.Query("SELECT DB_NAME()")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var rowValue string
		err = rows.Scan(&rowValue)
		if err != nil {
			break
		}
		log.Printf("SELECT DB_NAME() returns: %s\n", rowValue)
	}
}
