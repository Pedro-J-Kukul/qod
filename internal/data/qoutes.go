package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	// v.Check(q.Type == "inspire" || q.Type == "management" || q.Type == "sports" || q.Type == "life" || q.Type == "funny", "type", "must be one of the supported types")

	// Check that the quote is provided
	v.Check(q.Qoute != "", "quote", "must be provided")
	// Check that the quote is not too long (e.g., max 1000 characters)
	v.Check(len(q.Qoute) <= 1000, "quote", "must not be more than 1000 characters long")

	// Check that the author is provided
	v.Check(q.Author != "", "author", "must be provided")
	// Check that the author is not too long (e.g., max 255 characters)
	v.Check(len(q.Author) <= 255, "author", "must not be more than 255 characters long")
}

type QuoteModel struct {
	DB *sql.DB
}

func (q QuoteModel) Insert(quote *Qoute) error {
	query := `
		INSERT INTO quotes (type, quote, author)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
	`

	args := []any{quote.Type, quote.Qoute, quote.Author}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.ID, &quote.CreatedAt, &quote.Version)
}

func (q QuoteModel) Get(id int64) (*Qoute, error) {
	// Check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// the SQL query to be executed
	query := `
		SELECT id, type, quote, author, created_at, version
		FROM quotes
		WHERE id = $1
	`

	// Declare a variable to hold the data returned by the query
	var quote Qoute

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query using QueryRowContext, passing in the context, query, and id
	err := q.DB.QueryRowContext(ctx, query, id).Scan(
		&quote.ID,
		&quote.Type,
		&quote.Qoute,
		&quote.Author,
		&quote.CreatedAt,
		&quote.Version,
	)

	// Check for which error was returned by the query
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &quote, nil
}

// Update Method
func (q QuoteModel) Update(quote *Qoute) error {
	// sql query to update a record
	query := `
		UPDATE quotes
		SET type = $1, quote = $2, author = $3, version = version + 1
		WHERE id = $4
		RETURNING version
	`
	// args slice to hold the values for the placeholders in the query
	args := []any{
		quote.Type,
		quote.Qoute,
		quote.Author,
		quote.ID,
	}
	// create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)
}

// Delete Method
func (q QuoteModel) Delete(id int64) error {
	// check if the id is valid
	if id < 1 {
		return ErrRecordNotFound
	}

	// sql query to delete a record
	query := `
		DELETE FROM quotes
		WHERE id = $1
	`

	// create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query
	result, err := q.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (q QuoteModel) GetAll(qtype string, quote string, author string) ([]*Qoute, error) {
	// sql query to get all records
	query := `
		SELECT id, type, quote, author, created_at, version
		FROM quotes
		WHERE (to_tsvector('simple', type) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', quote) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple', author) @@ plainto_tsquery('simple', $3) OR $3 = '')
		ORDER BY id
		`
	// create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query
	rows, err := q.DB.QueryContext(ctx, query, qtype, quote, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// create a slice to hold the quotes
	var quotes []*Qoute

	// iterate over the rows
	for rows.Next() {
		var quote Qoute
		if err := rows.Scan(
			&quote.ID,
			&quote.Type,
			&quote.Qoute,
			&quote.Author,
			&quote.CreatedAt,
			&quote.Version,
		); err != nil {
			return nil, err
		}
		quotes = append(quotes, &quote)
	}

	// check for errors from iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quotes, nil
}
