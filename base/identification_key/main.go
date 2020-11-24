package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var SRC = fmt.Sprintf(
	"host=%v port=%v user=%v password=%v dbname=%v",
	"127.0.0.1",
	"5432",
	"frost",
	"frost",
	"test",
)

func init() {
	var err error
	DB, err = sql.Open("postgres", SRC)
	if err != nil {
		log.Fatalf("DB Open:%v\n", err)
	}
}

func main() {

}

func NextCustomerID() int64 {
	var id int64
	row := DB.QueryRow("SELECT nextval('customer_id_seq');")
	if err := row.Scan(&id); err != nil {
		log.Fatalf("")
	}
	return id
}

func New(id int64, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

type User struct {
	id   int64
	name string
}
