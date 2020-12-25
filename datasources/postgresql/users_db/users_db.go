package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	driverName = "postgres"
	datasourceName = "users_db"
)

var (
	Client *sql.DB
)
func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "postgres", "p0stgres", datasourceName)
	var err error
	Client, err = sql.Open(driverName, psqlInfo)
	if err != nil {
		panic(err)
	}
	if errConn := Client.Ping(); errConn != nil {
		panic(errConn)
	}
	log.Println("Database succesfully connected")
}