package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var SRS string = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s",
	"127.0.0.1",
	"5432",
	"frost",
	"frost",
	"test",
)

func init() {
	var err error
	DB, err = sql.Open("postgres", SRS)
	if err != nil {
		log.Fatalf("PG Open: %v", err)
	}
}

func main() {

	g := new(ContactsTableGateway)

	// if err := g.Delete(5); err != nil {
	// 	log.Fatalf("Gateway Delete: %v", err)
	// }

	// id, err := g.Insert("Сидор", "89095556677", "home")
	// if err != nil {
	// 	log.Fatalf("Gateway Insert: %v", err)
	// }
	// fmt.Printf("Id: %v\n", id)

	// rows, err := g.FindAll()
	// if err != nil {
	// 	log.Fatalf("Gateway FindAll: %v", err)
	// }

	// for _, row := range rows {
	// 	printRow(row)
	// }

	// row, err := g.Find(1)
	// if err != nil {
	// 	log.Fatalf("Gateway err: %v\n", err)
	// }
	// printRow(row)

	// rows, err := g.FindWhere("WHERE kind = 'home'")
	// if err != nil {
	// 	log.Fatalf("Gateway FindWhere: %v", err)
	// }
	// for _, r := range rows {
	// 	printRow(r)
	// }

	err := g.Update(1, "Марина", "89934443322", "home")
	if err != nil {
		log.Fatalf("Gateway Update: %v", err)
	}
	row, err := g.Find(1)
	if err != nil {
		log.Fatalf("Gateway Find: %v", err)
	}
	printRow(row)

}

func printRow(row Row) {
	fmt.Printf("%v %v %v %v\n", row["id"], row["name"], row["phone"], row["kind"])
}

type Row map[string]interface{}

type RowSet map[int]Row

type ContactsTableGateway struct{}

func (g *ContactsTableGateway) Update(id int, name, phone, kind string) error {
	_, err := DB.Exec(
		"UPDATE contacts"+
			" SET name=$1, phone=$2, kind=$3"+
			" WHERE id=$4",
		name, phone, kind, id)
	return err
}

func (g *ContactsTableGateway) FindWhere(whereClause string) (RowSet, error) {
	rows, err := DB.Query("SELECT id, name, phone, kind FROM contacts " + whereClause)
	if err != nil {
		return RowSet{}, err
	}
	result := RowSet{}
	for rows.Next() {
		var id int
		var name string
		var phone string
		var kind string
		if err := rows.Scan(&id, &name, &phone, &kind); err != nil {
			log.Fatalf("Find Where Scan: %v\n", err)
		}
		result[id] = Row{
			"id":    id,
			"name":  name,
			"phone": phone,
			"kind":  kind,
		}
	}
	return result, nil
}

func (g *ContactsTableGateway) Find(id int) (Row, error) {
	row := DB.QueryRow("SELECT id, name, phone, kind FROM contacts WHERE id = $1;", id)
	var rowid int
	var name string
	var phone string
	var kind string
	if err := row.Scan(&rowid, &name, &phone, &kind); err != nil {
		return Row{}, err
	}
	res := Row{
		"id":    rowid,
		"name":  name,
		"phone": phone,
		"kind":  kind,
	}
	return res, nil
}

func (g *ContactsTableGateway) Delete(id int) error {
	res, err := DB.Exec("DELETE FROM contacts WHERE id = $1;", id)
	log.Println(res)
	return err
}

func (g *ContactsTableGateway) Insert(name, phone, kind string) (int, error) {
	row := DB.QueryRow(
		"INSERT INTO contacts(name, phone, kind) VALUES($1, $2, $3) RETURNING id;",
		name, phone, kind,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (g *ContactsTableGateway) FindAll() (RowSet, error) {
	results := RowSet{}

	rows, err := DB.Query("SELECT id, name, phone, kind FROM contacts")
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var id int
		var name string
		var phone string
		var kind string
		if err := rows.Scan(&id, &name, &phone, &kind); err != nil {
			log.Fatalf("PG Scan: %v", err)
		}
		results[id] = Row{
			"id":    id,
			"name":  name,
			"phone": phone,
			"kind":  kind,
		}
	}
	return results, nil
}
