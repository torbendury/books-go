package storage

import "github.com/torbendury/books-go/data"

// Storage is a backend agnostic interface for any kind of data store which allows CRUD operations.
// This allows for decoupled logic between server and persistence.
type Storage interface {
	Create(*data.Book) error
	Get(int) (*data.Book, error)
	GetAll() []data.Book
	Update(*data.Book) (*data.Book, error)
	Delete(int) error
}
