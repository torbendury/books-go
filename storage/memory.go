// Package storage contains interfaces and implementations for different storage usage.
// At the time of writing, in-memory (temporary) and PostgreSQL is implemented.
package storage

import (
	"fmt"

	"github.com/torbendury/books-go/data"
)

// InMemoryStorage holds a in-memory slice which contains Books.
// Note: The InMemoryStorage is being thrown away when the application is stopped and therefore is not intended for any kind of usage beside testing.
type InMemoryStorage struct {
	Database []data.Book
	idSerial int
}

// NewInMemoryStorage returns a new InMemoryStorage pointer, initialized with an empty database which is represented by an empty slice of Books.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Database: make([]data.Book, 0),
		idSerial: 0,
	}
}

// Get iterates over the internal database and returns a book which matches the ID. If no book is found, an error is thrown.
func (ims *InMemoryStorage) Get(id int) (*data.Book, error) {
	for _, book := range ims.Database {
		if book.ID == id {
			return &book, nil
		}
	}
	return nil, fmt.Errorf("book id %v not found", id)
}

// GetAll returns the whole database of books.
func (ims *InMemoryStorage) GetAll() []data.Book {
	return ims.Database
}

// Create creates a new book in the InMemoryStorage.
// To implement the interface of a Storage, it is able to return an error.
func (ims *InMemoryStorage) Create(b *data.Book) (*data.Book, error) {
	ims.idSerial++
	b.ID = ims.idSerial
	ims.Database = append(ims.Database, *b)
	return b, nil
}

// Update checks if the given book exists by searching the database for its ID. If it is found, the entry in the database is replaced by the given Book.
func (ims *InMemoryStorage) Update(b *data.Book) (*data.Book, error) {
	if _, err := ims.Get(b.ID); err != nil {
		return nil, fmt.Errorf("book id %v not found", b.ID)
	}
	for idx, book := range ims.Database {
		if book.ID == b.ID {
			ims.Database[idx] = *b
			return b, nil
		}
	}
	return nil, fmt.Errorf("book id %v not found", b.ID)
}

// Delete checks if the given book exists by searching the database for its ID. If the ID is found, the book is being deleted.
// NOTE: The book is not actively being deleted. Rather than that, the last element of the slice is being put into the slice index where
// the to-be-deleted book resides. Then, the database is being cut down by the last element, effectively "deleting" the requested element.
// Since InMemoryStorages only purpose is for local testing, this is not an issue and allows for better DELETE performance.
func (ims *InMemoryStorage) Delete(id int) error {
	for idx, book := range ims.Database {
		if book.ID == id {
			ims.Database[idx] = ims.Database[len(ims.Database)-1]
			ims.Database = ims.Database[:len(ims.Database)-1]
			return nil
		}
	}
	return fmt.Errorf("book id %v not found", id)
}
