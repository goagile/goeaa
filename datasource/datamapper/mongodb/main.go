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
	ctx := context.Background()
	Client = NewMongoDBClient(ctx, uri)
	DB = Client.Database("testdatabase")
	Books = DB.Collection("books")
	defer Client.Disconnect(ctx)

	mapper := &BookDataMapper{}

	//
	// DeleteAll
	//
	if err := mapper.DeleteAll(ctx); err != nil {
		log.Fatal("DeleteAll", err)
	}

	//
	// Book Domain Entity
	//
	book := &Book{
		Title:     "Портрет Дориана Грея",
		Author:    "Уайльд Оскар",
		BasePrice: 1000,
	}

	//
	// Insert
	//
	if err := mapper.Insert(ctx, book); err != nil {
		log.Println("Insert", err)
	}
	fmt.Println("Insert", strJSON(book))

	//
	// Find By ID
	//
	foundByID, err := mapper.FindByID(ctx, book.ObjectID)
	if err != nil {
		log.Fatal("FindByID", err)
	}
	fmt.Println("foundByID", strJSON(foundByID))

	// Apply Some Business Logic
	foundByID.SetDiscountedPriceByPercent(50)

	//
	// Update
	//
	if err := mapper.UpdateDiscountedPrice(ctx, foundByID); err != nil {
		log.Fatal("UpdateDiscountedPrice", err)
	}

	//
	// Find By ID
	//
	updatedByID, err := mapper.FindByID(ctx, foundByID.ObjectID)
	if err != nil {
		log.Fatal("FindByID", err)
	}
	fmt.Println("updatedByID", strJSON(updatedByID))
}

func strJSON(book *Book) string {
	byt, err := json.MarshalIndent(book, "  ", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}
	return string(byt)
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

// BookDataMapper - ...
type BookDataMapper struct{}

// DeleteAll - ...
func (m *BookDataMapper) DeleteAll(ctx context.Context) error {
	_, err := Books.DeleteMany(ctx, bson.M{})
	return err
}

// Insert - ...
func (m *BookDataMapper) Insert(ctx context.Context, b *Book) error {
	r, err := Books.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	b.ObjectID = r.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID - ...
func (m *BookDataMapper) FindByID(ctx context.Context, id interface{}) (*Book, error) {
	r := Books.FindOne(ctx, bson.M{"_id": id})
	if err := r.Err(); err != nil {
		return nil, err
	}
	var b *Book
	if err := r.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateDiscountedPrice - ...
func (m *BookDataMapper) UpdateDiscountedPrice(ctx context.Context, b *Book) error {
	filter := bson.M{"_id": b.ObjectID}
	update := bson.M{"$set": bson.M{
		"discountedprice": b.DiscountedPrice,
	}}
	if _, err := Books.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// Book - ...
type Book struct {
	ObjectID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string             `json:"title" bson:"title,omitempty"`
	Author          string             `json:"author" bson:"author,omitempty"`
	BasePrice       int                `json:"baseprice" bson:"baseprice,omitempty"`
	DiscountedPrice int                `json:"discountedprice" bson:"discountedprice,omitempty"`
}

// SetDiscountedPriceByPercent - ...
func (b *Book) SetDiscountedPriceByPercent(p int) {
	if p > 100 {
		p = 100
	} else if p < 0 {
		p = 0
	}
	discount := float64(p) / 100.0
	b.DiscountedPrice = int(float64(b.BasePrice) - float64(b.BasePrice)*discount)
}
