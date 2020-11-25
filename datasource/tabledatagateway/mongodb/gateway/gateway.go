package gateway

import (
	"context"
	"log"

	"github.com/goagile/goeaa/datasource/tabledatagateway/mongodb/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// NewMongoDB - ...
func NewMongoDB(ctx context.Context, uri string) *gw {
	return NewMongoDBFromClient(NewMongoDBClient(ctx, uri))
}

// NewMongoDBFromClient - ...
func NewMongoDBFromClient(client *mongo.Client) *gw {
	db := client.Database("testdatabase")
	books := db.Collection("books")
	authors := db.Collection("authors")
	g := &gw{
		client:  client,
		db:      db,
		books:   books,
		authors: authors,
	}
	return g
}

type gw struct {
	client  *mongo.Client
	db      *mongo.Database
	books   *mongo.Collection
	authors *mongo.Collection
}

//
// Book methods
//
func (g *gw) DeleteAllBooks(ctx context.Context) (int64, error) {
	r, err := g.books.DeleteMany(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return r.DeletedCount, nil
}

func (g *gw) UpdateBookPrice(
	ctx context.Context,
	id interface{},
	base, discounted, discount int,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"price": bson.M{
			"base":       base,
			"discounted": discounted,
		},
		"discount": discount,
	}}
	_, err := g.books.UpdateOne(ctx, filter, update)
	return err
}

func (g *gw) UpdateAuthorBio(
	ctx context.Context,
	id interface{},
	bio string,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"bio": bio,
	}}
	_, err := g.authors.UpdateOne(ctx, filter, update)
	return err
}

func (g *gw) FindBookByISBN(ctx context.Context, isbn string) (*dto.Book, error) {
	r := g.books.FindOne(ctx, bson.M{"isbn": isbn})
	if err := r.Err(); err != nil {
		return new(dto.Book), err
	}
	var b *dto.Book
	if err := r.Decode(&b); err != nil {
		return new(dto.Book), err
	}
	return b, nil
}

func (g *gw) InsertBook(ctx context.Context, b *dto.Book) (interface{}, error) {
	r, err := g.books.InsertOne(ctx, b)
	if err != nil {
		return 0, err
	}
	return r.InsertedID, nil
}

func (g *gw) CountBooks(ctx context.Context) (int64, error) {
	c, err := g.books.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return c, nil
}

//
// Authors methods
//
func (g *gw) DeleteAllAuthors(ctx context.Context) (int64, error) {
	r, err := g.authors.DeleteMany(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return r.DeletedCount, nil
}

func (g *gw) InsertAuthor(ctx context.Context, a *dto.Author) (interface{}, error) {
	r, err := g.authors.InsertOne(ctx, a)
	if err != nil {
		return 0, err
	}
	return r.InsertedID, nil
}

func (g *gw) FindAuthorByID(ctx context.Context, id interface{}) (*dto.Author, error) {
	r := g.authors.FindOne(ctx, bson.M{"_id": id})
	if err := r.Err(); err != nil {
		return new(dto.Author), err
	}
	var a *dto.Author
	if err := r.Decode(&a); err != nil {
		return new(dto.Author), nil
	}
	return a, nil
}

func (g *gw) Disconnect(ctx context.Context) {
	g.client.Disconnect(ctx)
}
