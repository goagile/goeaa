package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

var DocMap = NewDocumentIdentityMap()

func init() {
	var err error
	DB, err = sql.Open("postgres", SRC)
	if err != nil {
		log.Fatalf("DB Open: %v\n", err)
	}
}

func main() {
	f := new(DocumentFinder)

	d, err := f.Find(1)
	if err != nil {
		log.Fatalf("DocumentFinder Find:%v\n", err)
	}
	d.Print()

	d, err = f.Find(1)
	if err != nil {
		log.Fatalf("DocumentFinder Find:%v\n", err)
	}
	d.Print()
}

type DocumentFinder struct{}

func (f *DocumentFinder) Find(id int) (*DocumentRowGateway, error) {
	if d, ok := DocMap.GetDocumentRowGateway(id); ok {
		fmt.Println("From IdMap")
		return d, nil
	}

	row := DB.QueryRow(
		"SELECT id, filename, createddate"+
			" FROM documents"+
			" WHERE id=$1;",
		id,
	)
	d := new(DocumentRowGateway)
	err := row.Scan(
		&d.ID,
		&d.Filename,
		&d.Createddate,
	)
	if err != nil {
		return d, err
	}

	DocMap.AddDocumentRowGateway(d)

	return d, nil
}

type DocumentRowGateway struct {
	ID          int
	Filename    string
	Createddate time.Time
}

func (d *DocumentRowGateway) Print() {
	fmt.Println("DocumentRowGateway:")
	fmt.Printf("\tID:%v\n", d.ID)
	fmt.Printf("\tFilename:%v\n", d.Filename)
	fmt.Printf("\tCreateddate:%v\n", d.Createddate.Format(time.RFC3339))
}

func NewDocumentIdentityMap() *DocumentIdentityMap {
	m := new(DocumentIdentityMap)
	m.data = map[int]*DocumentRowGateway{}
	return m
}

type DocumentIdentityMap struct {
	data map[int]*DocumentRowGateway
}

func (m *DocumentIdentityMap) GetDocumentRowGateway(id int) (*DocumentRowGateway, bool) {
	d, ok := m.data[id]
	return d, ok
}

func (m *DocumentIdentityMap) AddDocumentRowGateway(d *DocumentRowGateway) {
	m.data[d.ID] = d
}
