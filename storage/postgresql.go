package storage

import (
	"fmt"

	"github.com/torbendury/books-go/data"
)

type PostgresqlStorage struct {
	Database []data.Book
}

func NewPostgresqlStorage() *PostgresqlStorage {
	return &PostgresqlStorage{
		Database: make([]data.Book, 0),
	}
}

func (psql *PostgresqlStorage) Get(id int) (*data.Book, error) {
	for _, book := range psql.Database {
		if book.ID == id {
			return &book, nil
		}
	}
	return nil, fmt.Errorf("book id %v not found", id)
}

func (psql *PostgresqlStorage) GetAll() []data.Book {
	return psql.Database
}

func (psql *PostgresqlStorage) Create(b *data.Book) error {
	book, err := psql.Get(b.ID)
	if err != nil {
		psql.Database = append(psql.Database, *b)
		return nil
	}
	return fmt.Errorf("book id %v already exists", book.ID)
}

func (psql *PostgresqlStorage) Update(b *data.Book) (*data.Book, error) {
	return nil, nil
}

func (psql *PostgresqlStorage) Delete(id int) error {
	return nil
}
