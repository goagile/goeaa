package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	DB *database
)

func main() {
	ctx := context.Background()

	DB = NewDB()
	DB.InsertBook(ctx, &Book{
		ID:    12345,
		Title: "Read This Please",
		Reviews: []*Review{{
			BookID:      12345,
			Username:    "oldschool.reader",
			Title:       "Не очень",
			Text:        "Книга не очень понравилась, Скучно.",
			CreatedDate: time.Now(),
		}},
	})

	s := NewBookReviewService()
	var request *MakeReviewRequest
	webdata := []byte(`{
		"bookid":   12345,
		"reviewtitle":    "Хорошая книга",
		"reviewtext":     "Очень хорошая книга, порекомендую всем друзьям",
		"reviewusername": "book.liker"
	}`)
	if err := json.Unmarshal(webdata, &request); err != nil {
		log.Fatal("Unmarshal", err)
	}

	if err := s.MakeReviewForBook(ctx, request); err != nil {
		log.Fatal("MakeReviewForBook", err)
	}

	found := DB.FindBook(ctx, 12345)
	fmt.Println(found)
}

// NewBookReviewService - ...
func NewBookReviewService() *service {
	return new(service)
}

type service struct{}

// MakeReviewForBook - ...
func (s *service) MakeReviewForBook(ctx context.Context, r *MakeReviewRequest) error {
	book := DB.FindBook(ctx, r.BookID)
	review := &Review{
		BookID:      book.ID,
		Title:       r.ReviewTitle,
		Text:        r.ReviewText,
		Username:    r.ReviewUsername,
		CreatedDate: time.Now(),
	}
	book.MakeReview(review)
	return nil
}

// MakeReviewRequest - ...
type MakeReviewRequest struct {
	BookID         int64  `json:"bookid"`
	ReviewTitle    string `json:"reviewtitle"`
	ReviewText     string `json:"reviewtext"`
	ReviewUsername string `json:"reviewusername"`
}

// NewDB - ...
func NewDB() *database {
	db := new(database)
	db.m = make(map[int64]*Book)
	return db
}

type database struct {
	m map[int64]*Book
}

// InsertBook - ...
func (db *database) InsertBook(ctx context.Context, b *Book) {
	db.m[b.ID] = b
}

// FindBook - ...
func (db *database) FindBook(ctx context.Context, id int64) *Book {
	found := db.m[id]
	return found
}

// Book - ...
type Book struct {
	ID      int64
	Title   string
	Reviews []*Review
}

// MakeReview - ...
func (b *Book) MakeReview(r *Review) {
	b.Reviews = append(b.Reviews, r)
}

// String
func (b *Book) String() string {
	s := []string{}
	for _, r := range b.Reviews {
		s = append(s, r.String())
	}
	return fmt.Sprintf(
		"Book(\n\t%v\n\n%v\n)",
		b.Title,
		strings.Join(s, "\n\n"),
	)
}

// Review - ...
type Review struct {
	BookID      int64     `json:"bookid"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	Username    string    `json:"username"`
	CreatedDate time.Time `json:"createddate"`
}

func (r *Review) String() string {
	return fmt.Sprintf(
		"\t@%v\t%v\n\t%v\n\t%v",
		r.Username,
		r.Title,
		r.Text,
		r.CreatedDate.Format(time.RFC3339),
	)
}
