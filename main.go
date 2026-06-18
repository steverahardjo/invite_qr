package main

import (
	"database/sql"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}

func init_db(connString string) *sql.DB {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
