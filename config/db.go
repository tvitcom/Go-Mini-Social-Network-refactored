package config

import (
	"database/sql"
	"os"

	// for mysql
	_ "github.com/go-sql-driver/mysql"
)

// DB function
func DB() *sql.DB {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	_db := os.Getenv("DB")

	db, _ := sql.Open("mysql", user+":"+password+"@tcp("+host+":3306)/"+_db)
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
