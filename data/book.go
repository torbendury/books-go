// Package data contains data structures which are needed and used over the whole project.
// Changes in the data structures will probably cause API to break without also changing the logic itself.
package data

// Struct book is a public struct which described a single book and its JSON representation.
type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
