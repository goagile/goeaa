package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/goagile/goeaa/datasource/tabledatagateway/mongodb/dto"
	"github.com/goagile/goeaa/datasource/tabledatagateway/mongodb/gateway"
)

const (
	uri = "mongodb://127.0.0.1:27017"
)

var book = &dto.Book{
	Title:     "Портрет Дориана Грея",
	Authors:   []*dto.Author{{Name: "Уайльд Оскар"}},
	Price:     &dto.Price{Base: 200, Discounted: 160},
	Discount:  20,
	PubOffice: &dto.PubOffice{Name: "AСТ", Year: 2015},
	ISBN:      "978-5-17-099056-6",
	PageCount: 320,
}

func main() {
	//
	// Initialize Gateway
	//
	ctx := context.Background()
	g := gateway.NewMongoDB(ctx, uri)
	defer g.Disconnect(ctx)

	//
	// Delete All Books
	//
	if _, err := g.DeleteAllBooks(ctx); err != nil {
		log.Fatal("DeleteAllBooks", err)
	}

	//
	// Insert One Book
	//
	id, err := g.InsertBook(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("id", id)

	//
	// FindBook
	//
	found, err := g.FindBookByISBN(ctx, "978-5-17-099056-6")
	if err != nil {
		log.Fatal("FindBook", err)
	}
	b, _ := json.MarshalIndent(found, "  ", "  ")
	fmt.Println("found", string(b))

	//
	// Count All Books
	//
	c, err := g.CountBooks(ctx)
	if err != nil {
		log.Fatal("CountBooks", err)
	}
	fmt.Println("CountBooks", c)

	//
	// Update
	//
	base := 250
	discounted := 225
	discount := 10
	if err := g.UpdateBookPrice(ctx, id, base, discounted, discount); err != nil {
		log.Fatal("UpdateBookPrice", err)
	}

	//
	// FindBook
	//
	updated, err := g.FindBookByISBN(ctx, "978-5-17-099056-6")
	if err != nil {
		log.Fatal("FindBook", err)
	}
	u, _ := json.MarshalIndent(updated, "  ", "  ")
	fmt.Println("found", string(u))
}
