package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host = "localhost"
	port = 5432
	user = "application_user"
	password = "password"
	dbname = "company"
)

func init() {
	var psqlconn string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Println("Connection: " + psqlconn)
	db, err := sql.Open("postgres", psqlconn)
	fmt.Printf("db: %v", db)
	if err != nil {
		panic(err.Error())
	}

	DB = db
}
