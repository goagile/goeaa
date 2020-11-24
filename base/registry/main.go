package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var REG = new(Registry)

func init() {
	src := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v",
		"127.0.0.1",
		"5432",
		"frost",
		"frost",
		"test",
	)
	db, err := sql.Open("postgres", src)
	if err != nil {
		log.Fatalf("DB Open: %v\n", err)
	}
	REG.SetBD(db)

	f := new(BookFinder)
	REG.SetBookFinder(f)
}

func main() {
	f := REG.GetBookFinder()

	b, err := f.FindBook(123456)
	if err != nil {
		log.Fatalf("FindBook: %v\n", err)
	}

	b.Print()
}

//
// Book
//
type Book struct {
	ISBN int
	Name string
}

func (b *Book) Print() {
	fmt.Println("Book:")
	fmt.Printf("\tISBN:%v\n", b.ISBN)
	fmt.Printf("\tName:%v\n", b.Name)
}

//
// Registry
//
type Registry struct {
	db         *sql.DB
	bookFinder *BookFinder
}

func (r *Registry) SetBD(db *sql.DB) {
	r.db = db
}

func (r *Registry) GetDB() *sql.DB {
	return r.db
}

func (r *Registry) SetBookFinder(b *BookFinder) {
	r.bookFinder = b
}

func (r *Registry) GetBookFinder() *BookFinder {
	return r.bookFinder
}

//
// BookFinder
//
type BookFinder struct{}

func (f *BookFinder) FindBook(isbn int) (*Book, error) {
	db := REG.GetDB()
	row := db.QueryRow(
		"SELECT isbn, name"+
			" FROM books"+
			" WHERE isbn=$1;",
		isbn,
	)
	b := new(Book)
	if err := row.Scan(&b.ISBN, &b.Name); err != nil {
		return b, err
	}
	return b, nil
}
