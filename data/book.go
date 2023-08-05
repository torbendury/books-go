// Package data contains data structures which are needed and used over the whole project.
// Changes in the data structures will probably cause API to break without also changing the logic itself.
package data

import "errors"

// Struct book is a public struct which described a single book and its JSON representation.
type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (b *Book) Validate() error {
	if b.ID < 0 {
		return errors.New("id must not be less than 0")
	}
	if len(b.Title) < 1 {
		return errors.New("title must not be empty")
	}
	if len(b.Description) < 1 {
		return errors.New("description must not be empty")
	}
	if b.Price <= 0 {
		return errors.New("price must not be less than or equal to 0")
	}
	return nil
}
