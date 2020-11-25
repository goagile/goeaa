package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
	Books  *mongo.Collection
)

const (
	uri = "mongodb://127.0.0.1:27017"
)

func main() {
	//
	// Prepare Database
	//
	ctx := context.Background()
	Client = NewMongoDBClient(ctx, uri)
	DB = Client.Database("testdatabase")
	Books = DB.Collection("books")
	defer Client.Disconnect(ctx)

	//
	// Clear
	//
	if err := DeleteAllBooks(ctx); err != nil {
		log.Fatal("DeleteAllBooks", err)
	}

	//
	// Book
	//
	book := &BookActiveRecord{
		Title:     "Портрет Дориана Грея",
		Author:    "Уайльд Оскар",
		BasePrice: 1000,
	}

	//
	// Insert
	//
	if err := book.Insert(ctx); err != nil {
		log.Fatal("Insert", err)
	}
	insertedBookJSON, _ := json.MarshalIndent(book, "  ", "  ")
	fmt.Println("insertedBook", string(insertedBookJSON))

	//
	// Find
	//
	found, err := FindBookByID(ctx, book.ObjectID)
	if err != nil {
		log.Fatal("FindByID", err)
	}
	foundJSON, _ := json.MarshalIndent(found, "  ", "  ")
	fmt.Println("foundJSON", string(foundJSON))

	// Apply Some Business Logic
	found.SetDiscountedPriceByPercent(25)

	//
	// Update
	//
	if err := found.Update(ctx); err != nil {
		log.Fatal("Update", err)
	}

	//
	// Find
	//
	foundAfterUpdate, err := FindBookByID(ctx, book.ObjectID)
	if err != nil {
		log.Fatal("FindByID", err)
	}
	foundAfterUpdateJSON, _ := json.MarshalIndent(foundAfterUpdate, "  ", "  ")
	fmt.Println("foundAfterUpdateJSON", string(foundAfterUpdateJSON))
}

// BookActiveRecord - ...
type BookActiveRecord struct {
	ObjectID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string             `json:"title" bson:"title,omitempty"`
	Author          string             `json:"author" bson:"author,omitempty"`
	BasePrice       int                `json:"baseprice" bson:"baseprice,omitempty"`
	DiscountedPrice int                `json:"discountedprice" bson:"discountedprice,omitempty"`
}

// SetDiscountedPriceByPercent - ...
func (b *BookActiveRecord) SetDiscountedPriceByPercent(p int) {
	if p > 100 {
		p = 100
	} else if p < 0 {
		p = 0
	}
	discount := float64(p) / 100.0
	b.DiscountedPrice = int(float64(b.BasePrice) - float64(b.BasePrice)*discount)
}

// Insert - ...
func (b *BookActiveRecord) Insert(ctx context.Context) error {
	r, err := Books.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	b.ObjectID = r.InsertedID.(primitive.ObjectID)
	return nil
}

// Update - ...
func (b *BookActiveRecord) Update(ctx context.Context) error {
	filter := bson.M{"_id": b.ObjectID}
	_, err := Books.ReplaceOne(ctx, filter, b)
	if err != nil {
		return err
	}
	return nil
}

// NewMongoDBClient - ...
func NewMongoDBClient(ctx context.Context, uri string) *mongo.Client {
	opts := options.Client().ApplyURI(uri)

	c, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal("NewClient", err)
	}

	if err := c.Connect(ctx); err != nil {
		log.Fatal("Connect", err)
	}

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatal("Ping", err)
	}

	return c
}

// FindBookByID - ...
func FindBookByID(ctx context.Context, id interface{}) (*BookActiveRecord, error) {
	r := Books.FindOne(ctx, bson.M{"_id": id})
	if err := r.Err(); err != nil {
		return nil, err
	}
	var b *BookActiveRecord
	if err := r.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteAllBooks - ...
func DeleteAllBooks(ctx context.Context) error {
	_, err := Books.DeleteMany(ctx, bson.M{})
	return err
}
