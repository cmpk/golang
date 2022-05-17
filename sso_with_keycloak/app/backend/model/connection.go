package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Connection *sql.DB
}

// func CreateConnection() *Connection {
func CreateConnection() *sql.DB {
	connection := Connection{}
	connection.Open()
	log.Println("Opened DB")
	return connection.Connection
}

func (db *Connection) Open() {
	connectionString :=
		fmt.Sprintf("%s:%s@(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	var err error
	db.Connection, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
}

func (db *Connection) Close() {
	if db.Connection != nil {
		db.Connection.Close()
	}
	log.Println("Closed DB")
}
