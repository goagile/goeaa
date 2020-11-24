package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/lib/pq"
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
		log.Fatalf("DB Open^%v\n", err)
	}
}

func main() {
	// p := NewPostActiveRecord("сегодня хорошая погода")
	// if _, err := p.Insert(); err != nil {
	// 	log.Fatalf("Active Record Insert:%v\n", err)
	// }
	// printPostActiveRecord(p)

	// f := new(PostFinder)
	// p, err := f.Find(2)
	// if err != nil {
	// 	log.Fatalf("Finder Find:%v\n", err)
	// }
	// printPostActiveRecord(p)

	// f := new(PostFinder)
	// p, err := f.Find(5)
	// if err != nil {
	// 	log.Fatalf("Finder Find:%v\n", err)
	// }
	// fmt.Println("before update")
	// printPostActiveRecord(p)

	// p.Text = "новый текст"
	// if err := p.Update(); err != nil {
	// 	log.Fatalf("Active Record Update:%v\n", err)
	// }

	// p, err = f.Find(5)
	// if err != nil {
	// 	log.Fatalf("Finder Find:%v\n", err)
	// }
	// fmt.Println("after update")
	// printPostActiveRecord(p)

	p := NewPostActiveRecord("эта запись должна быть удалена")
	if _, err := p.Insert(); err != nil {
		log.Fatalf("Active Record Insert:%v\n", err)
	}
	printPostActiveRecord(p)
	if err := p.Delete(); err != nil {
		log.Fatalf("Active Record Delete:%v\n", err)
	}
}

type PostFinder struct{}

func (f *PostFinder) Find(id int) (*PostActiveRecord, error) {
	row := DB.QueryRow("SELECT id, text, createddate FROM posts WHERE id=$1;", id)
	p := new(PostActiveRecord)
	var createddate pq.NullTime
	if err := row.Scan(&p.ID, &p.Text, &createddate); err != nil {
		return p, err
	}
	p.Createddate = createddate.Time
	p.findWords()
	return p, nil
}

func printPostActiveRecord(p *PostActiveRecord) {
	fmt.Printf(
		"Post\n"+
			"\tID:%v\n"+
			"\tText:%q\n"+
			"\tAt:%v\n"+
			"\tWords:%v\n"+
			"\tCount:%v\n",
		p.ID,
		p.Text,
		p.Createddate.Format(time.RFC3339),
		p.Words,
		p.CountWords(),
	)
}

var re = regexp.MustCompile("[А-Яа-яA-Za-z]+")

func NewPostActiveRecord(text string) *PostActiveRecord {
	p := new(PostActiveRecord)
	p.ID = -1
	p.Text = text
	p.Createddate = time.Now()
	p.findWords()
	return p
}

type PostActiveRecord struct {
	ID          int
	Text        string
	Createddate time.Time
	Words       []string
}

func (p *PostActiveRecord) findWords() {
	p.Words = re.FindAllString(p.Text, -1)
}

func (p *PostActiveRecord) CountWords() int {
	return len(p.Words)
}

func (p *PostActiveRecord) Insert() (int, error) {
	row := DB.QueryRow(
		"INSERT INTO posts(text, createddate)"+
			" VALUES($1, $2)"+
			" RETURNING id;",
		p.Text,
		p.Createddate,
	)
	if err := row.Scan(&p.ID); err != nil {
		return p.ID, err
	}
	return p.ID, nil
}

func (p *PostActiveRecord) Update() error {
	_, err := DB.Exec(
		"UPDATE posts"+
			" SET text=$1"+
			" WHERE id=$2",
		p.Text,
		p.ID,
	)
	return err
}

func (p *PostActiveRecord) Delete() error {
	_, err := DB.Exec("DELETE FROM posts WHERE id=$1;", p.ID)
	return err
}
