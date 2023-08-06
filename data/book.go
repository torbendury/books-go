// Package data contains data structures which are needed and used over the whole project.
// Changes in the data structures will probably cause API to break without also changing the logic itself.
package data

// Struct book is a public struct which described a single book and its JSON representation.
type Book struct {
	ID          int     `json:"id" validate:"numeric,min=0"`
	Title       string  `json:"title" validate:"required,min=1"`
	Description string  `json:"description" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,numeric,min=0"`
}

type BookValidationError struct {
	Field string
	Tag   string
	Value string
}
