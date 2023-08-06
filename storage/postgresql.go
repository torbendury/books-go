package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/torbendury/books-go/data"
)

// PostgresqlStorage holds a connection to a PostgreSQL database instance.
type PostgresqlStorage struct {
	databaseConnection *sql.DB
}

// NewPostgresqlStorage returns a new PostgresqlStorage pointer, initialized with a database connection.
func NewPostgresqlStorage(db *sql.DB) *PostgresqlStorage {
	return &PostgresqlStorage{
		databaseConnection: db,
	}
}

// OpenDB takes connection information for a reachable PostgreSQL database, opens a connection to it and return it if the connection has been established.
func OpenDB(dbHost string, dbPort int, dbUser string, dbPass string, dbName string) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		// TODO: we might later try to recover from this, but right now we want to fail.
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		// TODO: again, we later might be able to escape from this by implementing retries?
		panic(err)
	}
	return db
}

// Get returns a book pointer if a matching book was found in the PSQL database. Otherwise, an error is raised.
func (psql *PostgresqlStorage) Get(id int) (*data.Book, error) {
	query := `
		SELECT id, title, description, price
		FROM books
		WHERE id = $1
	`
	var book data.Book
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := psql.databaseConnection.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.Price,
	)
	if err != nil {
		return nil, err
	}
	return &book, err
}

// GetAll returns all stored books from the PostgreSQL database.
// TODO: Implement limiting and pagination.
func (psql *PostgresqlStorage) GetAll() []data.Book {
	query := `
		SELECT id, title, description, price
		FROM books
	`
	books := make([]data.Book, 0)
	rows, err := psql.databaseConnection.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var book data.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Price); err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return books
}

// Create creates a new book in the PostgreSQL database and returns it, including its ID.
func (psql *PostgresqlStorage) Create(b *data.Book) (*data.Book, error) {
	query := `
		INSERT INTO books(title, description, price) 
		VALUES ($1, $2, $3)
		RETURNING id, title, description, price
	`
	var resultBook data.Book
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := psql.databaseConnection.QueryRowContext(ctx, query, b.Title, b.Description, b.Price).Scan(
		&resultBook.ID,
		&resultBook.Title,
		&resultBook.Description,
		&resultBook.Price,
	)
	if err != nil {
		return nil, err
	}
	return &resultBook, nil
}

// Update checks if the given book exists by its ID. If it is found, the entry is being updated. Otherwise an error is raised.
func (psql *PostgresqlStorage) Update(b *data.Book) (*data.Book, error) {
	query := `
		UPDATE books
		SET title = $2, description = $3, price = $4
		WHERE id = $1
		RETURNING id, title, description, price
	`
	var resultBook data.Book
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := psql.databaseConnection.QueryRowContext(ctx, query, b.ID, b.Title, b.Description, b.Price).Scan(
		&resultBook.ID,
		&resultBook.Title,
		&resultBook.Description,
		&resultBook.Price,
	)
	if err != nil {
		return nil, err
	}
	return &resultBook, nil
}

// Delete looks up a book in the PostgreSQL database and deletes it. If deletion fails, an error is returned.
// Also, if no rows are affected (i.e. because the book ID does not exist), an error is returned.
func (psql *PostgresqlStorage) Delete(id int) error {
	query := `
		DELETE FROM books
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	res, err := psql.databaseConnection.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
