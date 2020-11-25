package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Book - main aggregate
type Book struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title,omitempty"`
	Authors   []*BookAuthor      `bson:"autrhors,omitempty"`
	Price     *Price             `bson:"price,omitempty"`
	Discount  int                `bson:"discount,omitempty"`
	PubOffice *PubOffice         `bson:"puboffice,omitempty"`
	ISBN      string             `bson:"isbn,omitempty"`
	PageCount int                `bson:"pagecount,omitempty"`
}

type BookAuthor struct {
	Name string `bson:"name,omitempty"`
}

type Price struct {
	Base       int `bson:"base,omitempty"`
	Discounted int `bson:"discounted,omitempty"`
}

type PubOffice struct {
	Name string `bson:"name,omitempty"`
	Year int    `bson:"year,omitempty"`
}

// Author - author of book
type Author struct {
	ObjectID primitive.ObjectID `bson:"_id,omitempty"`
	First    string             `bson:"first,omitempty"`
	Last     string             `bson:"last,omitempty"`
	Bio      string             `bson:"bio,omitempty"`
}
