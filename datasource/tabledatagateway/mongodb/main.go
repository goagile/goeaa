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
	Authors:   []*dto.BookAuthor{{Name: "Уайльд Оскар"}},
	Price:     &dto.Price{Base: 200, Discounted: 160},
	Discount:  20,
	PubOffice: &dto.PubOffice{Name: "AСТ", Year: 2015},
	ISBN:      "978-5-17-099056-6",
	PageCount: 320,
}

var author = &dto.Author{
	First: "Оскар",
	Last:  "Уайльд",
	Bio:   "Оскар Уайльд родился в 1854 году в Дублине.",
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
	bookID, err := g.InsertBook(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bookID", bookID)

	//
	// FindBook
	//
	foundBook, err := g.FindBookByISBN(ctx, "978-5-17-099056-6")
	if err != nil {
		log.Fatal("FindBook", err)
	}
	foundBookJSON, _ := json.MarshalIndent(foundBook, "  ", "  ")
	fmt.Println("foundBookJSON", string(foundBookJSON))

	//
	// Count All Books
	//
	countBooks, err := g.CountBooks(ctx)
	if err != nil {
		log.Fatal("CountBooks", err)
	}
	fmt.Println("countBooks", countBooks)

	//
	// Update
	//
	base := 250
	discounted := 225
	discount := 10
	if err := g.UpdateBookPrice(ctx, bookID, base, discounted, discount); err != nil {
		log.Fatal("UpdateBookPrice", err)
	}

	//
	// FindBook
	//
	updatedBook, err := g.FindBookByISBN(ctx, "978-5-17-099056-6")
	if err != nil {
		log.Fatal("FindBookByISBN", err)
	}
	updatedBookJSON, _ := json.MarshalIndent(updatedBook, "  ", "  ")
	fmt.Println("updatedBookJSON", string(updatedBookJSON))

	//
	// Delete All Authors
	//
	if _, err := g.DeleteAllAutrhors(ctx); err != nil {
		log.Fatal("DeleteAllAutrhors", err)
	}

	//
	// Insert New Author
	//
	authorID, err := g.InsertAuthor(ctx, author)
	if err != nil {
		log.Fatal("InsertAuthor", err)
	}
	fmt.Println("authorID", authorID)

	//
	// Find Author
	//
	foundAuthor, err := g.FindAuthorByID(ctx, authorID)
	if err != nil {
		log.Fatal("FindAuthorByID", err)
	}
	foundAuthorJSON, _ := json.MarshalIndent(foundAuthor, "  ", "  ")
	fmt.Println("foundAuthorJSON", string(foundAuthorJSON))

	//
	// Update
	//
	newBio := "..."
	if err := g.UpdateAuthorBio(ctx, authorID, newBio); err != nil {
		log.Fatal("UpdateAuthorBio", err)
	}

	//
	// Find Author
	//
	udatedAuthor, err := g.FindAuthorByID(ctx, authorID)
	if err != nil {
		log.Fatal("FindAuthorByID", err)
	}
	udatedAuthorJSON, _ := json.MarshalIndent(udatedAuthor, "  ", "  ")
	fmt.Println("udatedAuthorJSON", string(udatedAuthorJSON))
}
