package driver

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("mysql", os.Getenv("DATA_SOURCE1"))
	fmt.Println(os.Getenv("DATA_SOURCE1"))
	if err != nil {
		fmt.Println("here", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
