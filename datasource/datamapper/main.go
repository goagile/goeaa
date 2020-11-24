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
	m := new(CustomerMapper)

	// c := new(Customer)
	// c.Name = "Iggos"
	// c.CountOrders = 18
	// printCustomer(c)

	// _, err := m.Insert(c)
	// if err != nil {
	// 	log.Printf("Mapper Insert:%v\n", err)
	// }
	// printCustomer(c)

	// c.ID = 4
	// c.CountOrders = 14
	// if err := m.Update(c); err != nil {
	// 	log.Printf("Mapper Update:%v\n", err)
	// }
	// printCustomer(c)

	// c.ID = 4
	// if err := m.Delete(c); err != nil {
	// 	log.Printf("Mapper Delete:%v\n", err)
	// }

	// c, err := m.Find(5)
	// if err != nil {
	// 	log.Printf("Mapper Find:%v\n", err)
	// }
	// printCustomer(c)

	cs, err := m.FindAll()
	if err != nil {
		log.Printf("Mapper FindAll:%v\n", err)
	}
	for _, c := range cs {
		printCustomer(c)
	}
}

type CustomerMapper struct{}

func (m *CustomerMapper) FindAll() (map[int]*Customer, error) {
	cs := map[int]*Customer{}
	rows, err := DB.Query(
		"SELECT id, name, count_orders" +
			" FROM customers;",
	)
	if err != nil {
		return cs, err
	}
	for rows.Next() {
		c := new(Customer)
		if err := rows.Scan(&c.ID, &c.Name, &c.CountOrders); err != nil {
			log.Printf("Scan: %v\n", err)
		}
		cs[c.ID] = c
	}
	return cs, nil
}

func (m *CustomerMapper) Find(id int) (*Customer, error) {
	row := DB.QueryRow(
		"SELECT name, count_orders"+
			" FROM customers"+
			" WHERE id=$1",
		id,
	)
	c := new(Customer)
	if err := row.Scan(&c.Name, &c.CountOrders); err != nil {
		return c, err
	}
	c.ID = id
	return c, nil
}

func (m *CustomerMapper) Delete(c *Customer) error {
	_, err := DB.Exec("DELETE FROM customers WHERE id=$1", c.ID)
	return err
}

func (m *CustomerMapper) Update(c *Customer) error {
	_, err := DB.Exec(
		"UPDATE customers"+
			" SET name=$1, count_orders=$2"+
			" WHERE id=$3",
		c.Name,
		c.CountOrders,
		c.ID,
	)
	return err
}

func (m *CustomerMapper) Insert(c *Customer) (int, error) {
	row := DB.QueryRow(
		"INSERT INTO customers(name, count_orders)"+
			" VALUES ($1, $2)"+
			" RETURNING id;",
		c.Name,
		c.CountOrders,
	)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	c.ID = id
	return id, nil
}

func printCustomer(c *Customer) {
	fmt.Printf(
		"Customer\n"+
			"\tID:%v\n"+
			"\tName:%v\n"+
			"\tCountOrders:%v\n"+
			"\tCountDiscount:%v\n",
		c.ID,
		c.Name,
		c.CountOrders,
		c.CountDiscount(),
	)
}

type Customer struct {
	ID          int
	Name        string
	CountOrders int
}

func (c *Customer) CountDiscount() float32 {
	if c.CountOrders > 10 {
		return 0.3
	}
	return 0
}
