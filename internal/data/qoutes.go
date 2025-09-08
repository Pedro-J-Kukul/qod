package data

import (
	"github.com/Pedro-J-Kukul/qod/internal/validator"
)

type Qoute struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Qoute     string `json:"quote"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	Version   int32  `json:"version"`
}

func ValidateQoute(v *validator.Validator, q *Qoute) {
	// Check that the type is provided
	v.Check(q.Type != "", "type", "must be provided")
	// Check that the type is one of the supported types
	v.Check(q.Type == "inspire" || q.Type == "management" || q.Type == "sports" || q.Type == "life" || q.Type == "funny", "type", "must be one of the supported types")

	// Check that the quote is provided
	v.Check(q.Qoute != "", "quote", "must be provided")
	// Check that the quote is not too long (e.g., max 1000 characters)
	v.Check(len(q.Qoute) <= 1000, "quote", "must not be more than 1000 characters long")

	// Check that the author is provided
	v.Check(q.Author != "", "author", "must be provided")
	// Check that the author is not too long (e.g., max 255 characters)
	v.Check(len(q.Author) <= 255, "author", "must not be more than 255 characters long")
}
