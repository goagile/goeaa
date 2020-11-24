package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/goagile/goeaa/datasource/tabledatagateway/mongodb/gateway"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	uri = "mongodb://127.0.0.1:27017"
)

var (
	wnp = bson.M{
		"title": "War and Peace",
		"authors": bson.A{
			bson.M{"name": "Leo Tolstoy"},
		},
		"price": bson.M{
			"base":       1000,
			"discounted": 900,
		},
		"discount": 10,
		"puboffice": bson.M{
			"name": "AST",
			"year": "2016",
		},
		"isbn":      "978-5-94074-819-9",
		"pagecount": 915,
	}
)

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
	id, err := g.InsertBook(ctx, wnp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("id", id)

	//
	// Count All Books
	//
	c, err := g.CountBooks(ctx)
	if err != nil {
		log.Fatal("CountBooks", err)
	}
	fmt.Println("CountBooks", c)

	//
	// FindBook
	//
	found, err := g.FindBook(ctx, bson.M{"title": wnp["title"]})
	if err != nil {
		log.Fatal("FindBook", err)
	}
	b, _ := json.MarshalIndent(found, "  ", "  ")
	fmt.Println("found", string(b))
}
