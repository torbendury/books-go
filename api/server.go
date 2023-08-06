// Package api contains the server logic. Its purpose is to handle Fiber logic and provide functionality to create new servers and start them up.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/torbendury/books-go/data"
	"github.com/torbendury/books-go/storage"
)

// Server holds information about any kind of data store which implements the storage.Storage interface.
// Also, it allows customizing the listenAddress - mainly to choose a proper port to listen on.
// Lastly, it holds a pointer to the Fiber app itself to act on it.
type Server struct {
	store         storage.Storage
	listenAddress string
	fiberApp      *fiber.App
	validator     *validator.Validate
}

// NewServer returns a new Server instance. This Server instance is basically idling until the Start() method is called.
func NewServer(store storage.Storage, listenAddress string, config fiber.Config) *Server {
	return &Server{
		store:         store,
		listenAddress: listenAddress,
		fiberApp:      fiber.New(config),
		validator:     validator.New(),
	}
}

// Start is responsible for configuring middleware, registering routes and putting the Fiber app in listen mode.
// Since this is effectively starting process forks etc. and is never really "ending", we can not unit-test this.
func (s *Server) Start() error {
	s.fiberApp.Use(
		logger.New(logger.Config{
			Format:        "{\"time\":${time}, \"latency\":\"${cust_latency}\", \"method\":\"${method}\", \"path\":\"${path}\", \"ip\":\"${ip}\", \"body\":${cust_reqbody}, \"useragent\":\"${ua}\", \"status\":${status}}\n",
			DisableColors: true,
			Output:        os.Stdout,
			CustomTags: map[string]logger.LogFunc{
				"cust_latency": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
					latency := data.Stop.Sub(data.Start)
					return output.WriteString(fmt.Sprintf("%v", latency))
				},
				"cust_reqbody": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
					body := c.Body()
					dst := &bytes.Buffer{}
					if err := json.Compact(dst, body); err != nil {
						return output.WriteString(string(body))
					}
					return output.WriteString(dst.String())
				},
			},
		}),
	)

	s.fiberApp.Post("/book", s.ValidateBook, s.handleCreateBook)
	s.fiberApp.Get("/book/:id", s.handleGetBookById)
	s.fiberApp.Get("/books", s.handleGetAllBooks)
	s.fiberApp.Put("/book", s.ValidateBook, s.handleUpdateBook)
	s.fiberApp.Delete("/book/:id", s.handleDeleteBook)

	return s.fiberApp.Listen(s.listenAddress)
}

// ValidateBook is a middleware handler for schema validation.
// This enables us to encapsulate user input validation without worrying about it in every handler.
// As seen in `Start()`, we can register it as a middleware handler easily.
func (s *Server) ValidateBook(c *fiber.Ctx) error {
	var errors []*data.BookValidationError
	body := new(data.Book)
	c.BodyParser(&body)
	err := s.validator.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el data.BookValidationError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}

// handleGetBookById checks if a correct ID has been requested, a book with the requested ID
// exists in the store and returns it if it exists.
func (s *Server) handleGetBookById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}
	book, err := s.store.Get(id)
	if err != nil {
		return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
	}
	return c.JSON(book)
}

// handleGetAllBooks calls the configured store and returns a JSON list of all existing books.
func (s *Server) handleGetAllBooks(c *fiber.Ctx) error {
	books := s.store.GetAll()
	return c.JSON(books)
}

// handleCreateBook validates the request body. If the body is not a valid book, an error is returned.
// If the request body is valid, it calls the store to persist the book.
// If any error occurs during persisting the book, the error is returned.
// If the book has been created, it is returned to the client.
func (s *Server) handleCreateBook(c *fiber.Ctx) error {
	book := new(data.Book)
	err := c.BodyParser(book)
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}
	resultBook, err := s.store.Create(book)
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}
	return c.Status(fiber.StatusAccepted).JSON(resultBook)
}

// handleDeleteBook validates the requested book ID. If it is valid, the store is called to check
// if there is a book with the given ID. If the ID is found, the book is deleted.
// If any error occurs, it is returned to the client.
func (s *Server) handleDeleteBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}
	if err = s.store.Delete(id); err != nil {
		return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

// handleUpdateBook validates the request body to be a book.
// If the body is valid, the book is being looked up in the store.
// If the book exists, it is updated and the updated book is returned to the client.
func (s *Server) handleUpdateBook(c *fiber.Ctx) error {
	book := new(data.Book)
	err := c.BodyParser(book)
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}
	book, err = s.store.Update(book)
	if err != nil {
		return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
	}
	return c.JSON(book)
}
