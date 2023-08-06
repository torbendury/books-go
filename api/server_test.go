package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/torbendury/books-go/data"
	"github.com/torbendury/books-go/storage"
)

var testCreateBook = data.Book{
	Title:       "Test1",
	Description: "Test1",
	Price:       1.11,
}

var testBookList = []data.Book{
	{
		Title:       "Test1",
		Description: "Test1",
		Price:       1.11,
		ID:          1,
	},
	{
		Title:       "Test1",
		Description: "Test1",
		Price:       1.11,
		ID:          2,
	},
}

var testBookId = 1

var testUpdateBook = data.Book{
	Title:       "Test2",
	Description: "Test2",
	Price:       2.22,
	ID:          1,
}

var invalidBook = `{"riesling": "schorle"}`

var nonExistingBook = data.Book{
	ID:          420,
	Title:       "Blazing it",
	Description: "Since 1997",
	Price:       42.0,
}

func setupServer() *Server {
	return NewServer(storage.NewInMemoryStorage(), ":3000", fiber.Config{})
}

func Test_handleCreateBook(t *testing.T) {
	// grab a fresh server
	server := setupServer()
	// register necessary route
	server.fiberApp.Post("/book", server.handleCreateBook)
	// do test request
	body, err := json.Marshal(testCreateBook)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("POST", "/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := server.fiberApp.Test(req, -1)
	assert.Equal(t, 202, resp.StatusCode)

	invalidBook := `{"riesling": "schorle"}`

	body, err = json.Marshal(invalidBook)
	if err != nil {
		t.Error(err)
	}
	req = httptest.NewRequest("POST", "/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 400, resp.StatusCode)
}

func Test_handleGetBookById(t *testing.T) {
	// grab a fresh server
	server := setupServer()
	// register necessary route
	server.fiberApp.Get("/book/:id", server.handleGetBookById)

	// insert test data
	_, err := server.store.Create(&testCreateBook)
	if err != nil {
		t.Error(err)
	}

	// correct request
	req := httptest.NewRequest("GET", fmt.Sprintf("%v/%d", "/book", testBookId), nil)
	resp, _ := server.fiberApp.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	var responseBook data.Book
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(body, &responseBook)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testCreateBook, responseBook)

	// no parameter, this should return 404 for "no route"
	req = httptest.NewRequest("GET", "/book/", nil)
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 404, resp.StatusCode)

	// no parameter, this should return 404
	req = httptest.NewRequest("GET", "/book/420", nil)
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 404, resp.StatusCode)

	// no parameter, this should return 404
	req = httptest.NewRequest("GET", "/book/schorle", nil)
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 400, resp.StatusCode)
}

func Test_handleGetAllBooks(t *testing.T) {
	// grab a fresh server
	server := setupServer()
	// register necessary route
	server.fiberApp.Get("/books", server.handleGetAllBooks)

	// insert test data
	_, err := server.store.Create(&testCreateBook)
	if err != nil {
		t.Error(err)
	}

	// insert test data
	_, err = server.store.Create(&testCreateBook)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("GET", "/books", nil)
	resp, _ := server.fiberApp.Test(req, -1)
	// check response data
	defer resp.Body.Close()
	var responseBook []data.Book
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(body, &responseBook)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testBookList, responseBook)
}

func Test_handleUpdateBook(t *testing.T) {
	// grab a fresh server
	server := setupServer()
	// register necessary route
	server.fiberApp.Put("/book", server.handleUpdateBook)

	// insert test data
	_, err := server.store.Create(&testCreateBook)
	if err != nil {
		t.Error(err)
	}

	// do test request
	body, err := json.Marshal(testUpdateBook)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("PUT", "/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := server.fiberApp.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)
	// check response data
	defer resp.Body.Close()
	var responseBook data.Book
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(body, &responseBook)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testUpdateBook, responseBook)

	// non existing book
	body, err = json.Marshal(nonExistingBook)
	if err != nil {
		t.Error(err)
	}
	req = httptest.NewRequest("PUT", "/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 404, resp.StatusCode)

	// invalid book
	body, err = json.Marshal(invalidBook)
	if err != nil {
		t.Error(err)
	}
	req = httptest.NewRequest("PUT", "/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 400, resp.StatusCode)
}

func Test_handleDeleteBook(t *testing.T) {
	// grab a fresh server
	server := setupServer()
	// register necessary route
	server.fiberApp.Delete("/book/:id", server.handleDeleteBook)

	// insert test data
	_, err := server.store.Create(&testCreateBook)
	if err != nil {
		t.Error(err)
	}

	// delete correct book
	req := httptest.NewRequest("DELETE", fmt.Sprintf("%v/%d", "/book", testBookId), nil)
	resp, _ := server.fiberApp.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)

	// nonexisting, this should return 404
	req = httptest.NewRequest("DELETE", "/book/420", nil)
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 404, resp.StatusCode)

	// nonsense, this should return 400
	req = httptest.NewRequest("DELETE", "/book/schorle", nil)
	resp, _ = server.fiberApp.Test(req, -1)
	assert.Equal(t, 400, resp.StatusCode)
}
