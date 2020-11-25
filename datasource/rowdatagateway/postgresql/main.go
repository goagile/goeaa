package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var DB *sql.DB

var CONN = fmt.Sprintf(
	"host=%v port=%v user=%v password=%v dbname=%v",
	"127.0.0.1",
	"5432",
	"frost",
	"frost",
	"test",
)

func init() {
	var err error
	DB, err = sql.Open("postgres", CONN)
	if err != nil {
		log.Fatalf("DB Open: %v", err)
	}
}

func main() {
	g := new(TaskGateway)
	g.Name = "TaskXXX"
	g.Username = "Frost"
	g.Createddate = time.Now()
	fmt.Println("Before Insert")
	printTaskGateway(g)
	_, err := g.Insert()
	if err != nil {
		log.Fatalf("Gateway Insert:%v\n", err)
	}
	fmt.Println("After Insert")
	printTaskGateway(g)

	f := new(TaskFinder)
	// g, err := f.Find(1)
	// if err != nil {
	// 	log.Fatalf("Finder Find:%v\n", err)
	// }
	// printTaskGateway(g)

	// g, _ := f.Find(3)
	// g.Name = "ZZ85"
	// g.Username = "Frost"
	// err := g.Update()
	// if err != nil {
	// 	log.Fatalf("Gateway Update:%v\n", err)
	// }

	// g, _ := f.Find(5)
	// if err := g.Delete(); err != nil {
	// 	log.Fatalf("Gateway Delete:%v\n", err)
	// }

	gs, err := f.FindAll()
	if err != nil {
		log.Fatalf("Finder FindAll:%v\n", err)
	}
	for _, g := range gs {
		printTaskGateway(g)
	}
}

type TaskFinder struct{}

func (f *TaskFinder) Find(id int) (*TaskGateway, error) {
	g := new(TaskGateway)
	row := DB.QueryRow("SELECT id, name, createddate, username"+
		" FROM tasks"+
		" WHERE id=$1;", id)
	var username sql.NullString
	var createddate pq.NullTime
	if err := row.Scan(
		&g.ID,
		&g.Name,
		&createddate,
		&username,
	); err != nil {
		return g, err
	}
	g.Username = username.String
	g.Createddate = createddate.Time
	return g, nil
}

func (g *TaskFinder) FindAll() (map[int]*TaskGateway, error) {
	results := map[int]*TaskGateway{}
	rows, err := DB.Query("SELECT id, name, createddate, username FROM tasks;")
	if err != nil {
		return results, err
	}
	for rows.Next() {
		g := new(TaskGateway)
		var createddate pq.NullTime
		var username sql.NullString
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&createddate,
			&username,
		); err != nil {
			log.Fatalf("Finder FindAll:%v\n", err)
		}
		g.Createddate = createddate.Time
		g.Username = username.String
		results[g.ID] = g
	}
	return results, nil
}

func printTaskGateway(g *TaskGateway) {
	fmt.Printf("%v %v %v %v\n",
		g.ID, g.Name, g.Createddate.Format(time.RFC3339), g.Username,
	)
}

type TaskGateway struct {
	ID          int
	Name        string
	Username    string
	Createddate time.Time
}

func (g *TaskGateway) Insert() (int, error) {
	row := DB.QueryRow("INSERT INTO tasks(name, createddate, username)"+
		" VALUES($1, $2, $3)"+
		" RETURNING id;",
		g.Name,
		g.Createddate,
		g.Username,
	)
	if err := row.Scan(&g.ID); err != nil {
		return g.ID, err
	}
	return g.ID, nil
}

func (g *TaskGateway) Update() error {
	if _, err := DB.Exec("UPDATE tasks"+
		" SET name=$1, createddate=$2, username=$3"+
		" WHERE id=$4",
		g.Name,
		g.Createddate,
		g.Username,
		g.ID,
	); err != nil {
		return err
	}
	return nil
}

func (g *TaskGateway) Delete() error {
	_, err := DB.Exec("DELETE FROM tasks WHERE id=$1;", g.ID)
	return err
}
